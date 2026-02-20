import { puzzleGames } from '../data/personal'
import Link from 'next/link'

export default function PuzzleGamesPage() {
  return (
    <main className="min-h-screen">
      <div className="container-main pt-20 pb-16">
        <Link href="/" className="text-xs text-text-tertiary dark:text-[#6e6e80] hover:text-text-secondary dark:hover:text-[#8e8ea0] transition-colors mb-8 inline-block">
          &larr; Back
        </Link>
        <h1 className="text-2xl font-semibold tracking-tight mb-8">Favorite Puzzle Games</h1>
        <ul className="space-y-2">
          {puzzleGames.map((game) => (
            <li key={game} className="text-sm text-text-secondary dark:text-[#a0a0b0]">
              {game}
            </li>
          ))}
        </ul>
      </div>
    </main>
  )
}
