// apps/web/src/lib/errorHelpers.ts

/** Common shape used by most API error objects */
export interface AppErrorLike {
  status?: number;
  message?: string;
  details?: unknown;
  cause?: unknown;
}

/** Base extractor: turns any error object into a readable string */
export function extractErrorMessage(err: unknown): string {
  const e = err as AppErrorLike | { message?: string };
  if (typeof e?.message === 'string' && e.message.trim()) return e.message.trim();

  // try to safely stringify for debugging
  try {
    return JSON.stringify(err);
  } catch {
    return 'Unexpected error';
  }
}

/** Alias for legacy imports (`extractErrorDetails`) */
export function extractErrorDetails(err: unknown): string {
  return extractErrorMessage(err);
}

/** Another alias for older imports (`getUserFriendlyError`) */
export function getUserFriendlyError(err: unknown): string {
  return extractErrorMessage(err);
}

/** Default export for `import extractErrorDetails from '@/lib/errorHelpers'` */
export default extractErrorMessage;
