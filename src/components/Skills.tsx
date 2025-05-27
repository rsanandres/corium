'use client'

import { motion } from 'framer-motion'

const skillCategories = [
  {
    title: 'Tools and Platforms',
    skills: ['Docker', 'Kubernetes', 'AWS SageMaker', 'GCP Vertex', 'Azure', 'Lambda', 'Jupyter', 'DataDog', 'SLURM']
  },
  {
    title: 'Programming Languages',
    skills: ['Python', 'SQL', 'R', 'C++', 'Go']
  },
  {
    title: 'Libraries and Frameworks',
    skills: ['TensorFlow', 'Keras', 'PyTorch (Lightning)', 'Ray', 'HuggingFace', 'Jupyter', 'LangChain', 'OpenAI API']
  }
]

export default function Skills() {
  return (
    <section className="bg-white">
      <div className="section-container">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
        >
          <h2 className="heading-2 mb-12 text-center">Technical Skills</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {skillCategories.map((category, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: index * 0.1 }}
                viewport={{ once: true }}
                className="bg-primary-50 rounded-lg p-6"
              >
                <h3 className="heading-3 mb-4 text-primary-900">{category.title}</h3>
                <div className="flex flex-wrap gap-2">
                  {category.skills.map((skill, idx) => (
                    <span
                      key={idx}
                      className="px-3 py-1 bg-white rounded-full text-primary-700 text-sm font-medium shadow-sm"
                    >
                      {skill}
                    </span>
                  ))}
                </div>
              </motion.div>
            ))}
          </div>
        </motion.div>
      </div>
    </section>
  )
} 