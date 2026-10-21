import { useState, useRef, useEffect, useMemo } from 'react'
import {
    Chart,
    LineController,
    LineElement,
    PointElement,
    CategoryScale,
    LinearScale,
    Tooltip,
    Legend,
} from 'chart.js'

Chart.register(LineController, LineElement, PointElement, CategoryScale, LinearScale, Tooltip, Legend)

type HeroData = Record<string, number>
type JsonRoot = Record<string, HeroData>

interface HeroTheme {
    color: string
    gradientStart: string
    gradientEnd: string
    avatarBg: string
    icon: string
}

const HERO_THEMES: Record<string, HeroTheme> = {
    'Barbarian King': {
        color: '#f59e0b',
        gradientStart: '#f59e0b',
        gradientEnd: '#d97706',
        avatarBg: 'rgba(245, 158, 11, 0.15)',
        icon: '👑',
    },
    'Archer Queen': {
        color: '#ec4899',
        gradientStart: '#ec4899',
        gradientEnd: '#db2777',
        avatarBg: 'rgba(236, 72, 153, 0.15)',
        icon: '🏹',
    },
    'Grand Warden': {
        color: '#38bdf8',
        gradientStart: '#38bdf8',
        gradientEnd: '#0284c7',
        avatarBg: 'rgba(56, 189, 248, 0.15)',
        icon: '🧙‍♂️',
    },
    'Royal Champion': {
        color: '#10b981',
        gradientStart: '#10b981',
        gradientEnd: '#059669',
        avatarBg: 'rgba(16, 185, 129, 0.15)',
        icon: '🛡️',
    },
    'Minion Prince': {
        color: '#a855f7',
        gradientStart: '#a855f7',
        gradientEnd: '#7e22ce',
        avatarBg: 'rgba(168, 85, 247, 0.15)',
        icon: '🦇',
    },
}

const LINE_PALETTE = [
    '#38bdf8', // sky blue
    '#ec4899', // pink
    '#f59e0b', // amber
    '#10b981', // emerald
    '#a855f7', // purple
    '#f43f5e', // rose
    '#06b6d4', // cyan
    '#84cc16', // lime
    '#eab308', // yellow
    '#6366f1', // indigo
    '#14b8a6', // teal
    '#fb923c', // orange
]

interface TrendsViewProps {
    structure: Record<string, Record<string, string[]>>
    dataFiles: Record<string, { default?: JsonRoot } | JsonRoot>
}

export function TrendsView({ structure, dataFiles }: TrendsViewProps) {
    const [selectedHero, setSelectedHero] = useState<string>('Archer Queen')
    const [dataset, setDataset] = useState<'meta100.json' | 'meta10k.json'>('meta10k.json')
    const [folder, setFolder] = useState<'meta' | 'metapairs'>('meta')
    const chartRef = useRef<HTMLCanvasElement | null>(null)
    const chartInstanceRef = useRef<Chart | null>(null)

    // Sorted chronological months (2025-01 -> 2026-10)
    const months = useMemo(() => {
        if (!structure[folder]) return []
        return Object.keys(structure[folder]).sort()
    }, [structure, folder])

    // Collect all distinct equipment names for this hero across all months
    const { items, monthlyStats } = useMemo(() => {
        const itemSet = new Set<string>()
        // month -> item -> { count, pct }
        const monthlyStats: Record<string, Record<string, { count: number; pct: number }>> = {}

        for (const m of months) {
            monthlyStats[m] = {}
            const key = `/public/${folder}/${m}/${dataset}`
            const mod = dataFiles[key]
            if (!mod) continue

            const json = ('default' in mod ? mod.default : mod) as JsonRoot | undefined
            if (!json || !json[selectedHero]) continue

            const heroData = json[selectedHero]
            const totalPicks = Object.values(heroData).reduce((acc, v) => acc + Number(v), 0)
            const playerCount = folder === 'metapairs' ? totalPicks : totalPicks / 2

            for (const [name, count] of Object.entries(heroData)) {
                itemSet.add(name)
                const num = Number(count)
                const pct = playerCount > 0 ? (num / playerCount) * 100 : 0
                monthlyStats[m][name] = { count: num, pct }
            }
        }

        return {
            items: Array.from(itemSet),
            monthlyStats,
        }
    }, [months, folder, dataset, selectedHero, dataFiles])

    // Rank items by recent month usage
    const rankedItems = useMemo(() => {
        const latestMonth = months.at(-1) || ''
        const latestData = monthlyStats[latestMonth] || {}
        return [...items].sort((a, b) => {
            const pctA = latestData[a]?.pct || 0
            const pctB = latestData[b]?.pct || 0
            return pctB - pctA
        })
    }, [items, monthlyStats, months])

    // Render chart
    useEffect(() => {
        if (!chartRef.current) return

        if (chartInstanceRef.current) {
            chartInstanceRef.current.destroy()
        }

        const ctx = chartRef.current.getContext('2d')
        if (!ctx) return

        const datasets = rankedItems.slice(0, 8).map((name, idx) => {
            const color = LINE_PALETTE[idx % LINE_PALETTE.length]
            const dataPoints = months.map((m) => {
                const stat = monthlyStats[m]?.[name]
                return stat ? Number(stat.pct.toFixed(1)) : 0
            })

            return {
                label: name,
                data: dataPoints,
                borderColor: color,
                backgroundColor: color,
                borderWidth: idx === 0 ? 3 : 2,
                pointRadius: 3,
                pointHoverRadius: 6,
                tension: 0.25,
            }
        })

        chartInstanceRef.current = new Chart(ctx, {
            type: 'line',
            data: {
                labels: months,
                datasets,
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                interaction: {
                    mode: 'index',
                    intersect: false,
                },
                plugins: {
                    legend: {
                        position: 'top',
                        labels: {
                            color: '#94a3b8',
                            font: { family: 'Inter, system-ui, sans-serif', size: 12 },
                            usePointStyle: true,
                            pointStyle: 'circle',
                            padding: 16,
                        },
                    },
                    tooltip: {
                        backgroundColor: 'rgba(15, 23, 42, 0.95)',
                        titleColor: '#f8fafc',
                        bodyColor: '#cbd5e1',
                        borderColor: '#334155',
                        borderWidth: 1,
                        padding: 12,
                        callbacks: {
                            label: (context) => ` ${context.dataset.label}: ${context.parsed.y}%`,
                        },
                    },
                },
                scales: {
                    x: {
                        grid: { color: 'rgba(51, 65, 85, 0.35)' },
                        ticks: {
                            color: '#94a3b8',
                            font: { family: 'Inter, system-ui, sans-serif', size: 11 },
                            maxRotation: 45,
                        },
                    },
                    y: {
                        min: 0,
                        max: 100,
                        grid: { color: 'rgba(51, 65, 85, 0.35)' },
                        ticks: {
                            color: '#94a3b8',
                            font: { family: 'Inter, system-ui, sans-serif', size: 11 },
                            callback: (v) => `${v}%`,
                        },
                    },
                },
            },
        })

        return () => {
            if (chartInstanceRef.current) {
                chartInstanceRef.current.destroy()
            }
        }
    }, [months, rankedItems, monthlyStats])

    const theme = HERO_THEMES[selectedHero] || HERO_THEMES['Barbarian King']

    return (
        <div className="trends-view-container">
            {/* Hero selection pill row */}
            <div className="trends-controls">
                <div className="trends-hero-pills">
                    {Object.keys(HERO_THEMES).map((hero) => {
                        const hTheme = HERO_THEMES[hero]
                        const isActive = selectedHero === hero
                        return (
                            <button
                                key={hero}
                                className={`trend-pill ${isActive ? 'active' : ''}`}
                                style={
                                    isActive
                                        ? {
                                              borderColor: hTheme.color,
                                              boxShadow: `0 0 16px ${hTheme.color}33`,
                                              color: '#fff',
                                          }
                                        : {}
                                }
                                onClick={() => setSelectedHero(hero)}
                            >
                                <span className="pill-icon">{hTheme.icon}</span>
                                <span className="pill-name">{hero}</span>
                            </button>
                        )
                    })}
                </div>

                <div className="trends-sub-controls">
                    <div className="selector-group">
                        <label>Category</label>
                        <select
                            value={folder}
                            onChange={(e) => setFolder(e.target.value as 'meta' | 'metapairs')}
                        >
                            <option value="meta">Single Equipment</option>
                            <option value="metapairs">Equipment Pairs</option>
                        </select>
                    </div>

                    <div className="selector-group">
                        <label>Scope</label>
                        <select
                            value={dataset}
                            onChange={(e) => setDataset(e.target.value as 'meta100.json' | 'meta10k.json')}
                        >
                            <option value="meta10k.json">Top 10,000 Players</option>
                            <option value="meta100.json">Top 100 Elite</option>
                        </select>
                    </div>
                </div>
            </div>

            {/* Main Trend Line Chart */}
            <div className="trend-chart-card">
                <div className="trend-chart-header">
                    <div className="trend-hero-badge">
                        <span className="trend-icon">{theme.icon}</span>
                        <div>
                            <h3>{selectedHero} Meta Evolution</h3>
                            <p className="trend-subtitle">
                                Tracking Pick Rate % across all 22 seasons ({months[0]} → {months.at(-1)})
                            </p>
                        </div>
                    </div>
                </div>

                <div className="trend-canvas-wrapper" style={{ height: '420px', position: 'relative' }}>
                    <canvas ref={chartRef}></canvas>
                </div>
            </div>

            {/* Performance Summary Cards */}
            <div className="trend-breakdown-grid">
                {rankedItems.slice(0, 6).map((name, idx) => {
                    const firstMonth = months[0] || ''
                    const lastMonth = months.at(-1) || ''
                    const firstPct = monthlyStats[firstMonth]?.[name]?.pct || 0
                    const lastPct = monthlyStats[lastMonth]?.[name]?.pct || 0
                    const diff = lastPct - firstPct
                    const color = LINE_PALETTE[idx % LINE_PALETTE.length]

                    return (
                        <div key={name} className="trend-stat-card">
                            <div className="trend-card-top">
                                <span className="trend-color-dot" style={{ backgroundColor: color }}></span>
                                <span className="trend-card-title">{name}</span>
                                <span className="trend-current-pct">{lastPct.toFixed(1)}%</span>
                            </div>
                            <div className="trend-card-sub">
                                <span>Recent Season</span>
                                <span className={diff >= 0 ? 'trend-up' : 'trend-down'}>
                                    {diff >= 0 ? `+${diff.toFixed(1)}%` : `${diff.toFixed(1)}%`} all-time
                                </span>
                            </div>
                        </div>
                    )
                })}
            </div>
        </div>
    )
}
