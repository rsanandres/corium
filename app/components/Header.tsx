'use client'

import { useState } from 'react'
import Link from 'next/link'
import { navItems } from '../data/navigation'
import ThemeToggle from './ThemeToggle'

export default function Header() {
  const [mobileOpen, setMobileOpen] = useState(false)

  const openContact = () => {
    window.dispatchEvent(new Event('openContactModal'))
    setMobileOpen(false)
  }

  return (
    <header className="sticky top-0 z-50 border-b border-border dark:border-border-dark bg-surface/80 dark:bg-surface-dark/80 backdrop-blur-md">
      <div className="container-main flex items-center justify-between h-14">
        <Link href="/" className="text-sm font-medium tracking-tight hover:text-accent transition-colors">
          Raphael San Andres
        </Link>

        {/* Desktop nav */}
        <nav className="hidden md:flex items-center gap-6">
          {navItems.map((item) => (
            <a
              key={item.href}
              href={item.href}
              className="text-sm text-text-secondary dark:text-[#8e8ea0] hover:text-text-primary dark:hover:text-[#e5e5ea] transition-colors"
            >
              {item.label}
            </a>
          ))}
          <button
            onClick={openContact}
            className="text-sm text-text-secondary dark:text-[#8e8ea0] hover:text-text-primary dark:hover:text-[#e5e5ea] transition-colors"
          >
            Contact
          </button>
          <ThemeToggle />
        </nav>

        {/* Mobile hamburger */}
        <div className="flex items-center gap-3 md:hidden">
          <ThemeToggle />
          <button
            onClick={() => setMobileOpen(!mobileOpen)}
            className="p-1.5 text-text-secondary dark:text-[#8e8ea0]"
            aria-label="Toggle menu"
          >
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="1.5">
              {mobileOpen ? (
                <path d="M5 5l10 10M15 5L5 15" />
              ) : (
                <path d="M3 6h14M3 10h14M3 14h14" />
              )}
            </svg>
          </button>
        </div>
      </div>

      {/* Mobile menu */}
      {mobileOpen && (
        <nav className="md:hidden border-t border-border dark:border-border-dark bg-surface dark:bg-surface-dark px-6 py-4 space-y-3">
          {navItems.map((item) => (
            <a
              key={item.href}
              href={item.href}
              onClick={() => setMobileOpen(false)}
              className="block text-sm text-text-secondary dark:text-[#8e8ea0] hover:text-text-primary dark:hover:text-[#e5e5ea]"
            >
              {item.label}
            </a>
          ))}
          <button
            onClick={openContact}
            className="block text-sm text-text-secondary dark:text-[#8e8ea0] hover:text-text-primary dark:hover:text-[#e5e5ea]"
          >
            Contact
          </button>
        </nav>
      )}
    </header>
  )
}
