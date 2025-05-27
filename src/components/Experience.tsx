'use client'

import { motion } from 'framer-motion'

const experiences = [
  {
    title: 'Founding ML Solutions Engineer & Technical Account Manager',
    company: 'Undisclosed Company',
    period: 'February 2024 - Present',
    description: [
      'Grew and led a customer support-facing Machine Learning Engineering team',
      'Managed Kubernetes Clusters and educated AI Startups on ML workloads',
      'Developed documentation and orchestrated meetings',
      'Developed a Kubernetes node repair program',
      'Created and maintained Kubernetes Toolings'
    ]
  },
  {
    title: 'Machine Learning Support Engineer',
    company: 'Weights and Biases',
    period: 'January 2023 - January 2024',
    description: [
      'Debugged and solved issues from ML Practitioners',
      'Triaged and traced bugs',
      'Managed customer requests',
      'Collaborated with cross-functional teams'
    ]
  }
]

export default function Experience() {
  return (
    <section className="bg-white">
      <div className="section-container">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
        >
          <h2 className="heading-2 mb-12 text-center">Professional Experience</h2>
          
          <div className="space-y-12">
            {experiences.map((exp, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: index * 0.1 }}
                viewport={{ once: true }}
                className="border-l-4 border-primary-200 pl-6"
              >
                <h3 className="heading-3 text-primary-900">{exp.title}</h3>
                <p className="text-xl text-primary-600 mt-1">{exp.company}</p>
                <p className="text-primary-500 mt-1">{exp.period}</p>
                <ul className="mt-4 space-y-2">
                  {exp.description.map((item, idx) => (
                    <li key={idx} className="body-text flex items-start">
                      <span className="h-1.5 w-1.5 rounded-full bg-primary-400 mt-2 mr-3 flex-shrink-0" />
                      {item}
                    </li>
                  ))}
                </ul>
              </motion.div>
            ))}
          </div>
        </motion.div>
      </div>
    </section>
  )
} 