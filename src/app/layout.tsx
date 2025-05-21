import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Raphael San Andres - Machine Learning Engineer',
  description: 'Personal website of Raphael San Andres, showcasing ML engineering expertise and professional achievements.',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@100..900&family=JetBrains+Mono:wght@100..800&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className={`${inter.className} bg-primary-50 text-primary-900`}>
        <main className="min-h-screen">
          {children}
        </main>
      </body>
    </html>
  )
} 