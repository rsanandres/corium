import { ProjectItem } from '../data/content'

interface ProjectListProps {
  items: ProjectItem[]
}

const GitHubIcon = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z" />
  </svg>
)

export default function ProjectList({ items }: ProjectListProps) {
  return (
    <div>
      {items.map((item) => (
        <div key={item.name} className="py-3.5 first:pt-0 last:pb-0 group">
          <div className="flex flex-col sm:flex-row sm:items-baseline sm:justify-between gap-1 mb-1.5">
            <div className="flex items-center gap-2">
              {item.url ? (
                <a
                  href={item.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm font-medium group-hover:text-accent transition-colors flex items-center gap-1.5"
                >
                  <span className="inline-block w-1.5 h-1.5 rounded-full bg-emerald-500 shrink-0" />
                  {item.name}
                </a>
              ) : (
                <p className="text-sm font-medium">{item.name}</p>
              )}
              {item.github && (
                <a
                  href={item.github}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-text-tertiary dark:text-[#6e6e80] hover:text-text-primary dark:hover:text-[#e5e5ea] transition-colors"
                  aria-label={`${item.name} GitHub repository`}
                >
                  <GitHubIcon />
                </a>
              )}
            </div>
            <p className="text-xs text-text-tertiary dark:text-[#6e6e80] font-mono shrink-0">{item.period}</p>
          </div>
          <p className="text-sm text-text-secondary dark:text-[#a0a0b0] leading-relaxed">{item.description}</p>
          {item.tech && (
            <p className="text-xs text-text-tertiary dark:text-[#6e6e80] mt-1.5">{item.tech}</p>
          )}
        </div>
      ))}
    </div>
  )
}
