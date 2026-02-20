'use client'

import React, { useEffect, useState } from 'react'

// Contact form is disabled in this environment (no email/CAPTCHA credentials).
// To re-enable, restore the original form with ReCAPTCHA and fetch('/api/contact').

const ContactModal: React.FC = () => {
  const [open, setOpen] = useState(false)

  useEffect(() => {
    const openModal = () => setOpen(true)
    window.addEventListener('openContactModal', openModal)
    return () => window.removeEventListener('openContactModal', openModal)
  }, [])

  const closeModal = () => setOpen(false)

  if (!open) return null

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50" onClick={closeModal}>
      <div
        className="bg-surface dark:bg-surface-dark border border-border dark:border-border-dark rounded-lg shadow-xl p-8 w-full max-w-md relative"
        onClick={(e) => e.stopPropagation()}
      >
        <button
          className="absolute top-3 right-3 text-text-tertiary dark:text-[#6e6e80] hover:text-text-primary dark:hover:text-[#e5e5ea] text-xl transition-colors"
          onClick={closeModal}
          aria-label="Close"
        >
          &times;
        </button>
        <h2 className="text-lg font-semibold mb-5">Contact Me</h2>
        <p className="text-sm text-text-secondary dark:text-[#a0a0b0]">
          The contact form is not available in this environment. Please reach out via LinkedIn or GitHub instead.
        </p>
      </div>
    </div>
  )
}

export default ContactModal
