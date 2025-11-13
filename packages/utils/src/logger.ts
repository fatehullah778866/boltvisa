type Level = 'DEBUG' | 'INFO' | 'WARN' | 'ERROR'

const ts = () => new Date().toISOString()

// Make errors & arrays readable
const normalize = (items: unknown[]) =>
  items.map((x) => {
    if (x instanceof Error) return { name: x.name, message: x.message, stack: x.stack }
    if (Array.isArray(x)) return { array: x }
    return x
  })

function log(level: Level, msg: string, ...meta: unknown[]) {
  const line = `[${ts()}] [${level}] ${msg}`
  const payload = normalize(meta)
  
  if (level === 'ERROR') console.error(line, ...payload)
  else if (level === 'WARN') console.warn(line, ...payload)
  else if (level === 'INFO') console.info(line, ...payload)
  else console.debug(line, ...payload)
}

export const logger = {
  debug: (m: string, ...a: unknown[]) => log('DEBUG', m, ...a),
  info:  (m: string, ...a: unknown[]) => log('INFO',  m, ...a),
  warn:  (m: string, ...a: unknown[]) => log('WARN',  m, ...a),
  error: (m: string, ...a: unknown[]) => log('ERROR', m, ...a),
}

export type Logger = typeof logger

