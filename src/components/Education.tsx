import { motion } from 'framer-motion'

const education = [
  {
    degree: 'Masters in Artificial Intelligence',
    school: 'Penn State',
    location: 'Remote',
    period: 'Expected June 2025',
    details: []
  },
  {
    degree: 'Bachelor of Science in Statistics',
    school: 'UCLA',
    location: 'Los Angeles, CA',
    period: 'June 2022',
    details: []
  }
]

export default function Education() {
  return (
    <section className="bg-primary-50">
      <div className="section-container">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
        >
          <h2 className="heading-2 mb-12 text-center">Education</h2>
          
          <div className="space-y-8">
            {education.map((edu, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: index * 0.1 }}
                viewport={{ once: true }}
                className="bg-white rounded-lg p-6 shadow-sm"
              >
                <h3 className="heading-3 text-primary-900">{edu.degree}</h3>
                <p className="text-xl text-primary-600 mt-1">{edu.school}</p>
                <div className="flex justify-between items-center mt-2 text-primary-500">
                  <span>{edu.location}</span>
                  <span>{edu.period}</span>
                </div>
                {edu.details.length > 0 && (
                  <ul className="mt-4 space-y-2">
                    {edu.details.map((detail, idx) => (
                      <li key={idx} className="body-text flex items-start">
                        <span className="h-1.5 w-1.5 rounded-full bg-primary-400 mt-2 mr-3 flex-shrink-0" />
                        {detail}
                      </li>
                    ))}
                  </ul>
                )}
              </motion.div>
            ))}
          </div>
        </motion.div>
      </div>
    </section>
  )
} 