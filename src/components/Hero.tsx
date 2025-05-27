'use client'

import { motion } from 'framer-motion'
import { FaGithub, FaLinkedin } from 'react-icons/fa'

export default function Hero() {
  return (
    <section className="relative overflow-hidden bg-primary-50">
      <div className="section-container">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          className="text-center"
        >
          <h1 className="heading-1 mb-6">
            Raphael San Andres
            <span className="block text-2xl sm:text-3xl font-normal mt-2 text-primary-600">
              Machine Learning Engineer
            </span>
          </h1>
          
          <p className="body-text max-w-2xl mx-auto mb-8">
            Highly skilled and experienced Machine Learning Engineer with a passion for developing 
            and deploying innovative ML solutions. Proven ability to support early-stage startups 
            and solve complex ML problems.
          </p>
          
          <div className="flex justify-center space-x-6">
            <a
              href="https://github.com/rsanandres"
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary-700 hover:text-primary-900 transition-colors"
            >
              <FaGithub className="w-8 h-8" />
            </a>
            <a
              href="https://www.linkedin.com/in/raphael-san-andres/"
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary-700 hover:text-primary-900 transition-colors"
            >
              <FaLinkedin className="w-8 h-8" />
            </a>
          </div>
        </motion.div>
      </div>
    </section>
  )
} 