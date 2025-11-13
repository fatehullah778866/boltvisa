import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import { ErrorBoundaryWrapper } from '../components/ErrorBoundaryWrapper'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Visa Help Center',
  description: 'Comprehensive visa application management system',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <ErrorBoundaryWrapper>{children}</ErrorBoundaryWrapper>
      </body>
    </html>
  )
}
