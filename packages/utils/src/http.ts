export async function safeJson<T>(res: Response): Promise<T> {
  const text = await res.text()
  try { 
    return JSON.parse(text) as T
  } catch { 
    throw new Error(`Invalid JSON: ${text.slice(0, 200)}`)
  }
}

export async function getJson<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, init)
  if (!res.ok) throw new Error(`HTTP ${res.status} ${res.statusText}`)
  const data = await safeJson<T>(res)
  return data
}

