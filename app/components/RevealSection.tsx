'use client'

import { useEffect, useRef } from 'react'

interface RevealSectionProps {
  children: React.ReactNode
  className?: string
  id?: string
}

export default function RevealSection({ children, className = '', id }: RevealSectionProps) {
  const ref = useRef<HTMLElement>(null)

  useEffect(() => {
    const el = ref.current
    if (!el) return

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          el.classList.add('visible')
          observer.unobserve(el)
        }
      },
      { threshold: 0.15 }
    )

    observer.observe(el)
    return () => observer.disconnect()
  }, [])

  return (
    <section ref={ref} id={id} className={`reveal ${className}`}>
      {children}
    </section>
  )
}
