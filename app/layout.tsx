import type { Metadata } from 'next'
import { GeistSans } from 'geist/font/sans'
import { GeistMono } from 'geist/font/mono'
import './globals.css'
import Script from 'next/script'
import Header from './components/Header'
import Footer from './components/Footer'
import ContactModal from './components/ContactModal'

export const metadata: Metadata = {
  title: 'Raphael San Andres - Machine Learning Engineer',
  description: 'Personal website of Raphael San Andres, a Machine Learning Engineer specializing in ML solutions and AI development.',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" className={`${GeistSans.variable} ${GeistMono.variable} scroll-smooth`}>
      <head>
        <Script
          src="https://www.googletagmanager.com/gtag/js?id=G-8ERJT22BXP"
          strategy="afterInteractive"
        />
        <Script id="google-analytics" strategy="afterInteractive">
          {`
            window.dataLayer = window.dataLayer || [];
            function gtag(){dataLayer.push(arguments);}
            gtag('js', new Date());
            gtag('config', 'G-8ERJT22BXP');
          `}
        </Script>
      </head>
      <body className="font-sans bg-surface dark:bg-surface-dark text-text-primary dark:text-[#e5e5ea] transition-colors duration-200">
        <Header />
        <ContactModal />
        {children}
        <Footer />
      </body>
    </html>
  )
}
