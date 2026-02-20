import { socialLinks } from '../data/content'
import { personalPages } from '../data/navigation'
import Link from 'next/link'

export default function Footer() {
  return (
    <footer className="border-t border-border dark:border-border-dark">
      <div className="container-main py-8 flex flex-col sm:flex-row items-center justify-between gap-4">
        <p className="text-sm text-text-tertiary dark:text-[#6e6e80]">
          Raphael San Andres
        </p>
        <div className="flex items-center gap-5">
          {personalPages.map((page) => (
            <Link
              key={page.href}
              href={page.href}
              className="text-xs text-text-tertiary dark:text-[#6e6e80] hover:text-text-secondary dark:hover:text-[#8e8ea0] transition-colors"
            >
              {page.label}
            </Link>
          ))}
          <span className="w-px h-3 bg-border dark:bg-border-dark" />
          {socialLinks.map((link) => (
            <a
              key={link.label}
              href={link.href}
              target="_blank"
              rel="noopener noreferrer"
              className="text-xs text-text-tertiary dark:text-[#6e6e80] hover:text-text-secondary dark:hover:text-[#8e8ea0] transition-colors"
            >
              {link.label}
            </a>
          ))}
        </div>
      </div>
    </footer>
  )
}
