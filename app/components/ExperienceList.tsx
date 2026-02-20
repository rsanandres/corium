import { ExperienceItem } from '../data/content'

interface ExperienceListProps {
  items: ExperienceItem[]
}

export default function ExperienceList({ items }: ExperienceListProps) {
  return (
    <div>
      {items.map((item) => (
        <div key={item.company} className="py-4 first:pt-0 last:pb-0">
          <div className="mb-2">
            <h3 className="text-base font-medium">{item.company}</h3>
            {item.subtitle && (
              <p className="text-xs text-text-tertiary dark:text-[#6e6e80]">{item.subtitle}</p>
            )}
          </div>

          {item.positions ? (
            <div className="space-y-4">
              {item.positions.map((pos) => (
                <div key={pos.title}>
                  <div className="flex flex-col sm:flex-row sm:items-baseline sm:justify-between gap-1 mb-2">
                    <p className="text-sm font-medium text-text-secondary dark:text-[#8e8ea0]">{pos.title}</p>
                    <p className="text-xs text-text-tertiary dark:text-[#6e6e80] font-mono">{pos.period}</p>
                  </div>
                  <ul className="space-y-1.5">
                    {pos.details.map((detail, i) => (
                      <li key={i} className="text-sm text-text-secondary dark:text-[#a0a0b0] leading-relaxed pl-4 relative before:content-[''] before:absolute before:left-0 before:top-[9px] before:w-1.5 before:h-1.5 before:rounded-full before:bg-border dark:before:bg-border-dark">
                        {detail}
                      </li>
                    ))}
                  </ul>
                </div>
              ))}
            </div>
          ) : (
            <>
              <div className="flex flex-col sm:flex-row sm:items-baseline sm:justify-between gap-1 mb-2">
                <p className="text-sm font-medium text-text-secondary dark:text-[#8e8ea0]">{item.title}</p>
                <p className="text-xs text-text-tertiary dark:text-[#6e6e80] font-mono">{item.period}</p>
              </div>
              {item.details && (
                <ul className="space-y-1.5">
                  {item.details.map((detail, i) => (
                    <li key={i} className="text-sm text-text-secondary dark:text-[#a0a0b0] leading-relaxed pl-4 relative before:content-[''] before:absolute before:left-0 before:top-[9px] before:w-1.5 before:h-1.5 before:rounded-full before:bg-border dark:before:bg-border-dark">
                      {detail}
                    </li>
                  ))}
                </ul>
              )}
            </>
          )}
        </div>
      ))}
    </div>
  )
}
