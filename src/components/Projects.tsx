import { motion } from 'framer-motion'

const projects = [
  {
    title: 'Bird Classifier and GAN',
    period: 'April 2023 - August 2023',
    description: 'Created a Bird Classifier with high accuracy and developed a DCGAN.',
    tags: ['Computer Vision', 'Deep Learning', 'GAN']
  },
  {
    title: 'Q-Learning in Custom OpenAI Gym',
    period: 'April 2023 - August 2023',
    description: 'Designed a custom OpenAI Gym Environment and implemented agent classes for reinforcement learning.',
    tags: ['Reinforcement Learning', 'OpenAI Gym', 'Q-Learning']
  },
  {
    title: 'Scientific Paper Summarization Tool',
    period: 'September 2023 - Present',
    description: 'Collected and parsed scientific paper PDFs, currently finetuning a Transformer Model for summarization.',
    tags: ['NLP', 'Transformers', 'PDF Processing']
  }
]

export default function Projects() {
  return (
    <section className="bg-primary-50">
      <div className="section-container">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
        >
          <h2 className="heading-2 mb-12 text-center">Projects</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {projects.map((project, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: index * 0.1 }}
                viewport={{ once: true }}
                className="bg-white rounded-lg p-6 shadow-sm"
              >
                <h3 className="heading-3 text-primary-900">{project.title}</h3>
                <p className="text-primary-500 mt-1">{project.period}</p>
                <p className="body-text mt-3">{project.description}</p>
                <div className="flex flex-wrap gap-2 mt-4">
                  {project.tags.map((tag, idx) => (
                    <span
                      key={idx}
                      className="px-2 py-1 bg-primary-100 rounded-md text-primary-700 text-sm font-medium"
                    >
                      {tag}
                    </span>
                  ))}
                </div>
              </motion.div>
            ))}
          </div>
          
          <div className="mt-12 text-center">
            <h3 className="heading-3 mb-4">Publications</h3>
            <p className="body-text">RL with Mario - Article</p>
            
            <h3 className="heading-3 mt-8 mb-4">Certifications</h3>
            <p className="body-text">CKA (Certified Kubernetes Administrator) - In Progress (Est. February 2025)</p>
          </div>
        </motion.div>
      </div>
    </section>
  )
} 