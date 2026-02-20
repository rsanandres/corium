import Hero from './components/Hero'
import RevealSection from './components/RevealSection'
import ExperienceList from './components/ExperienceList'
import ProjectList from './components/ProjectList'
import { about, experience, skillCategories, education, projects, publications } from './data/content'

export default function Home() {
  return (
    <main className="min-h-screen">
      <Hero />

      {/* About */}
      <RevealSection id="about" className="container-main pb-6">
        <h2 className="text-lg font-semibold pb-2 mb-3 border-b border-border dark:border-border-dark">About</h2>
        <p className="text-sm text-text-secondary dark:text-[#a0a0b0] leading-relaxed">
          {about}
        </p>
      </RevealSection>

      {/* Experience */}
      <RevealSection id="experience" className="container-main py-6">
        <h2 className="text-lg font-semibold pb-2 mb-3 border-b border-border dark:border-border-dark">Experience</h2>
        <ExperienceList items={experience} />
      </RevealSection>

      {/* Projects */}
      <RevealSection id="projects" className="container-main py-6">
        <h2 className="text-lg font-semibold pb-2 mb-3 border-b border-border dark:border-border-dark">Projects</h2>
        <ProjectList items={projects} />
      </RevealSection>

      {/* Skills */}
      <RevealSection id="skills" className="container-main py-6">
        <h2 className="text-lg font-semibold pb-2 mb-3 border-b border-border dark:border-border-dark">Skills</h2>
        <div className="space-y-3">
          {skillCategories.map((category) => (
            <div key={category.name} className="flex flex-wrap items-baseline gap-x-2 gap-y-1.5">
              <span className="text-xs font-medium text-text-tertiary dark:text-[#6e6e80] mr-1 w-28 shrink-0">
                {category.name}
              </span>
              {category.skills.map((skill) => (
                <span
                  key={skill}
                  className="px-2 py-0.5 rounded text-xs text-text-secondary dark:text-[#a0a0b0] bg-black/[0.04] dark:bg-white/[0.06]"
                >
                  {skill}
                </span>
              ))}
            </div>
          ))}
        </div>
      </RevealSection>

      {/* Education */}
      <RevealSection className="container-main py-6">
        <h2 className="text-lg font-semibold pb-2 mb-3 border-b border-border dark:border-border-dark">Education</h2>
        <div className="space-y-3">
          {education.map((item) => (
            <div key={item.degree} className="flex flex-col sm:flex-row sm:items-baseline sm:justify-between gap-1">
              <p className="text-sm">{item.degree}</p>
              <p className="text-xs text-text-tertiary dark:text-[#6e6e80] font-mono">
                {item.school} &middot; {item.year}
              </p>
            </div>
          ))}
        </div>
      </RevealSection>

      {/* Publications */}
      <RevealSection className="container-main pb-12">
        <h2 className="text-lg font-semibold pb-2 mb-3 border-b border-border dark:border-border-dark">Publications</h2>
        {publications.map((pub) => (
          <a
            key={pub.title}
            href={pub.url}
            target="_blank"
            rel="noopener noreferrer"
            className="text-sm text-text-secondary dark:text-[#a0a0b0] hover:text-accent transition-colors"
          >
            {pub.title} ({pub.type}) &rarr;
          </a>
        ))}
      </RevealSection>
    </main>
  )
}
