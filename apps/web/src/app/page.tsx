import Link from 'next/link'

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <div className="z-10 max-w-5xl w-full items-center justify-between font-mono text-sm">
        <h1 className="text-4xl font-bold text-center mb-8">
          🧭 Visa Help Center
        </h1>
        <p className="text-center text-lg mb-4">
          Welcome to the Visa Application Management System
        </p>
        <div className="mt-8 text-center">
          <Link
            href="/login"
            className="px-6 py-3 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors inline-block"
          >
            Get Started
          </Link>
        </div>
      </div>
    </main>
  )
}

