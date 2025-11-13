// packages/utils/src/apiRequest.ts
export type HttpMethod = "GET" | "POST" | "PUT" | "DELETE";

export type AppError = {
  status: number;
  message: string;
  details?: unknown;
  cause?: unknown;
};

const BASE =
  (typeof process !== "undefined" && (process.env.NEXT_PUBLIC_API_URL || process.env.API_BASE)) ||
  (typeof window !== "undefined" ? (window as any).__API_BASE__ : "") ||
  "";

if (!BASE) {
  throw new Error(
    'NEXT_PUBLIC_API_URL is not set. Add it to apps/web/.env.local, e.g.\nNEXT_PUBLIC_API_URL=http://127.0.0.1:8080'
  );
}

export async function apiRequest<T>(
  path: string,
  opts: {
    method?: HttpMethod;
    token?: string;
    body?: unknown;
    headers?: Record<string, string>;
    timeoutMs?: number;
  } = {}
): Promise<T> {
  const {
    method = "GET",
    token,
    body,
    headers = {},
    timeoutMs = 15000, // 15s to avoid dev cold-start blips
  } = opts;

  const url = path.startsWith("http") ? path : `${BASE}${path}`;
  const controller = new AbortController();
  const t = setTimeout(() => controller.abort(new Error("timeout")), timeoutMs);

  const h: Record<string, string> = {
    "Content-Type": "application/json",
    ...headers,
  };
  if (token) h["Authorization"] = `Bearer ${token}`;

  let res: Response;

  try {
    res = await fetch(url, {
      method,
      headers: h,
      body: body !== undefined ? JSON.stringify(body) : undefined,
      signal: controller.signal,
    });
  } catch (err: any) {
    clearTimeout(t);
    const isTimeout =
      err?.name === "AbortError" ||
      err?.message?.toLowerCase?.().includes("timeout") ||
      err?.message === "timeout";
    const appError: AppError = {
      status: 0,
      message: isTimeout
        ? "Request timed out. Is the backend running on NEXT_PUBLIC_API_URL?"
        : err?.message || "Network error contacting API",
      details: { url, timeoutMs, method },
      cause: err,
    };
    throw appError;
  } finally {
    clearTimeout(t);
  }

  const text = await res.text();
  let json: any = {};
  try {
    json = text ? JSON.parse(text) : {};
  } catch {
    json = { message: text };
  }

  if (!res.ok) {
    const msg = json?.error || json?.message || `HTTP ${res.status}`;
    const appError: AppError = {
      status: res.status,
      message: msg,
      details: json,
    };
    throw appError;
  }

  return (json?.data ?? json) as T;
}
