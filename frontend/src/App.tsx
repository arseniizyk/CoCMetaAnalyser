import { useEffect, useState, useRef, useMemo } from 'react'
import './App.css'
import { TrendsView } from './TrendsView'
import {
    Chart,
    BarController,
    BarElement,
    CategoryScale,
    LinearScale,
    Tooltip,
} from 'chart.js'

Chart.register(BarController, BarElement, CategoryScale, LinearScale, Tooltip)

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

const DEFAULT_THEME: HeroTheme = {
    color: '#6366f1',
    gradientStart: '#6366f1',
    gradientEnd: '#4f46e5',
    avatarBg: 'rgba(99, 102, 241, 0.15)',
    icon: '⚡',
}

// Read bundled data files
const dataFiles = import.meta.glob('/public/*/*/*.json', { eager: true }) as Record<
    string,
    { default?: JsonRoot } | JsonRoot
>

const structure: Record<string, Record<string, string[]>> = {}

for (const path of Object.keys(dataFiles)) {
    const parts = path.split('/')
    const folder = parts.at(-3)!
    const month = parts.at(-2)!
    const file = parts.at(-1)!
    if (!structure[folder]) structure[folder] = {}
    if (!structure[folder][month]) structure[folder][month] = []
    if (!structure[folder][month].includes(file)) {
        structure[folder][month].push(file)
    }
}

// Sort months descending (latest season first)
for (const folder of Object.keys(structure)) {
    for (const m of Object.keys(structure[folder])) {
        structure[folder][m].sort()
    }
}

function formatFileName(fileName: string): string {
    if (fileName.includes('10k')) return 'Top 10,000 Players'
    if (fileName.includes('100')) return 'Top 100 Players'
    return fileName.replace('.json', '')
}

function HeroCard({
    hero,
    items,
    searchQuery,
}: {
    hero: string
    items: HeroData
    searchQuery: string
}) {
    const [viewMode, setViewMode] = useState<'chart' | 'list'>('chart')
    const chartRef = useRef<HTMLCanvasElement | null>(null)
    const chartInstanceRef = useRef<Chart | null>(null)

    const theme = HERO_THEMES[hero] || DEFAULT_THEME

    // Filter items by search and sort by pick count descending
    const sortedItems = useMemo(() => {
        return Object.entries(items)
            .filter(([name]) =>
                name.toLowerCase().includes(searchQuery.toLowerCase())
            )
            .sort((a, b) => Number(b[1]) - Number(a[1]))
    }, [items, searchQuery])

    const totalPicks = useMemo(() => {
        return Object.values(items).reduce((acc, val) => acc + Number(val), 0)
    }, [items])

    const topPick = sortedItems[0]

    useEffect(() => {
        if (viewMode !== 'chart' || !chartRef.current) return

        if (chartInstanceRef.current) {
            chartInstanceRef.current.destroy()
        }

        const labels = sortedItems.map(([name]) => name)
        const dataValues = sortedItems.map(([, count]) => Number(count))

        const ctx = chartRef.current.getContext('2d')
        if (!ctx) return

        chartInstanceRef.current = new Chart(ctx, {
            type: 'bar',
            data: {
                labels,
                datasets: [
                    {
                        label: hero,
                        data: dataValues,
                        backgroundColor: theme.color,
                        borderRadius: 6,
                        maxBarThickness: 18,
                    },
                ],
            },
            options: {
                indexAxis: 'y',
                responsive: true,
                maintainAspectRatio: false,
                animation: {
                    duration: 400,
                },
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        backgroundColor: '#1e293b',
                        titleColor: '#f8fafc',
                        bodyColor: '#94a3b8',
                        borderColor: 'rgba(255, 255, 255, 0.1)',
                        borderWidth: 1,
                        padding: 10,
                        callbacks: {
                            label: (context) => {
                                const val = Number(context.raw)
                                const pct = totalPicks > 0 ? ((val / totalPicks) * 100).toFixed(1) : '0'
                                return ` Usage: ${val.toLocaleString()} (${pct}%)`
                            },
                        },
                    },
                },
                scales: {
                    x: {
                        grid: {
                            color: 'rgba(255, 255, 255, 0.05)',
                        },
                        ticks: {
                            color: '#64748b',
                            font: { size: 11 },
                        },
                    },
                    y: {
                        grid: { display: false },
                        ticks: {
                            color: '#cbd5e1',
                            font: { size: 11, weight: 'bold' },
                        },
                    },
                },
            },
        })

        return () => {
            if (chartInstanceRef.current) {
                chartInstanceRef.current.destroy()
                chartInstanceRef.current = null
            }
        }
    }, [sortedItems, theme.color, hero, totalPicks, viewMode])

    return (
        <div className="hero-card" style={{ borderColor: `${theme.color}30` }}>
            <div className="hero-card-header">
                <div className="hero-info">
                    <div
                        className="hero-avatar"
                        style={{
                            background: theme.avatarBg,
                            border: `1px solid ${theme.color}40`,
                        }}
                    >
                        {theme.icon}
                    </div>
                    <div className="hero-name-wrap">
                        <h2>{hero}</h2>
                        <p className="hero-stats-sub">
                            {topPick
                                ? `Meta #1: ${topPick[0]} (${totalPicks > 0 ? ((Number(topPick[1]) / totalPicks) * 100).toFixed(1) : 0}%)`
                                : 'No equipment'}
                        </p>
                    </div>
                </div>

                <div className="card-view-toggle">
                    <button
                        className={`view-toggle-btn ${viewMode === 'chart' ? 'active' : ''}`}
                        onClick={() => setViewMode('chart')}
                        title="Chart view"
                    >
                        Chart
                    </button>
                    <button
                        className={`view-toggle-btn ${viewMode === 'list' ? 'active' : ''}`}
                        onClick={() => setViewMode('list')}
                        title="Leaderboard view"
                    >
                        List
                    </button>
                </div>
            </div>

            {sortedItems.length === 0 ? (
                <div className="state-container">No matching equipment found.</div>
            ) : viewMode === 'chart' ? (
                <div className="chart-container">
                    <canvas ref={chartRef}></canvas>
                </div>
            ) : (
                <div className="leaderboard-list">
                    {sortedItems.map(([name, count], idx) => {
                        const num = Number(count)
                        const pct = totalPicks > 0 ? (num / totalPicks) * 100 : 0
                        return (
                            <div key={name} className="leaderboard-row">
                                <div className="row-top">
                                    <span className="equip-label">
                                        <span className={`rank-badge rank-${idx + 1}`}>
                                            {idx === 0 ? '👑' : `#${idx + 1}`}
                                        </span>
                                        {name}
                                    </span>
                                    <div className="equip-counts">
                                        <span className="pick-count">{num.toLocaleString()}</span>
                                        <span className="pick-percent">{pct.toFixed(1)}%</span>
                                    </div>
                                </div>
                                <div className="progress-track">
                                    <div
                                        className="progress-bar-fill"
                                        style={{
                                            width: `${Math.min(pct, 100)}%`,
                                            background: `linear-gradient(to right, ${theme.gradientStart}, ${theme.gradientEnd})`,
                                        }}
                                    ></div>
                                </div>
                            </div>
                        )
                    })}
                </div>
            )}
        </div>
    )
}

export default function App() {
    const [activeTab, setActiveTab] = useState<'snapshot' | 'trends'>('snapshot')
    const folders = Object.keys(structure)
    const [folder, setFolder] = useState<string>(folders[0] || 'meta')

    const months = useMemo(() => {
        if (!structure[folder]) return []
        return Object.keys(structure[folder]).sort().reverse()
    }, [folder])

    const [month, setMonth] = useState<string>(months[0] || '')

    const files = useMemo(() => {
        if (!structure[folder] || !structure[folder][month]) return []
        return structure[folder][month]
    }, [folder, month])

    const [fileName, setFileName] = useState<string>(files[0] || '')
    const [data, setData] = useState<JsonRoot | null>(null)
    const [error, setError] = useState<string | null>(null)
    const [selectedHero, setSelectedHero] = useState<string>('All')
    const [searchQuery, setSearchQuery] = useState<string>('')

    // Keep month and file in sync when folder changes
    useEffect(() => {
        if (months.length > 0 && !months.includes(month)) {
            setMonth(months[0])
        }
    }, [months, month])

    useEffect(() => {
        if (files.length > 0 && !files.includes(fileName)) {
            setFileName(files[0])
        }
    }, [files, fileName])

    // Load data from memory or fallback fetch
    useEffect(() => {
        if (!folder || !month || !fileName) return

        setError(null)
        const eagerKey = `/public/${folder}/${month}/${fileName}`
        const eagerModule = dataFiles[eagerKey]

        if (eagerModule) {
            const rawData = 'default' in eagerModule ? eagerModule.default : eagerModule
            if (rawData) {
                setData(rawData as JsonRoot)
                return
            }
        }

        // Fallback to network fetch if not in eager glob
        const path = `${import.meta.env.BASE_URL}${folder}/${month}/${fileName}`
        fetch(path)
            .then((res) => {
                if (!res.ok) throw new Error(`HTTP error ${res.status}`)
                return res.json()
            })
            .then((json) => setData(json as JsonRoot))
            .catch((err) => setError(String(err)))
    }, [folder, month, fileName])

    const heroNames = useMemo(() => {
        if (!data) return []
        return Object.keys(data)
    }, [data])

    const filteredHeroes = useMemo(() => {
        if (!data) return []
        const entries = Object.entries(data)
        if (selectedHero === 'All') return entries
        return entries.filter(([hero]) => hero === selectedHero)
    }, [data, selectedHero])

    return (
        <div className="app-container">
            {/* Header */}
            <header className="app-header">
                <div className="brand-section">
                    <div className="logo-badge">⚔️</div>
                    <div className="brand-titles">
                        <h1>CoC Meta Analyser</h1>
                        <p>Clash of Clans Legendary League Hero Equipment Analytics</p>
                    </div>
                </div>

                <div className="header-meta">
                    <div className="season-badge">
                        <span>📅</span> Season {month || '—'}
                    </div>
                    <a
                        href="https://github.com/arseniizyk/CoCMetaAnalyser"
                        target="_blank"
                        rel="noreferrer"
                        className="github-btn"
                    >
                        ⭐ GitHub
                    </a>
                </div>
            </header>

            {/* Top View Mode Tabs */}
            <div className="app-tabs-bar">
                <button
                    className={`tab-btn ${activeTab === 'snapshot' ? 'active' : ''}`}
                    onClick={() => setActiveTab('snapshot')}
                >
                    <span>📊</span> Monthly Snapshot
                </button>
                <button
                    className={`tab-btn ${activeTab === 'trends' ? 'active' : ''}`}
                    onClick={() => setActiveTab('trends')}
                >
                    <span>📈</span> Historical Trends (2025–2026)
                </button>
            </div>

            {activeTab === 'trends' ? (
                <TrendsView structure={structure} dataFiles={dataFiles} />
            ) : (
                <>
                    {/* Controls Panel */}
                    <div className="controls-panel">
                <div className="control-group">
                    {/* Mode Tabs: meta vs metapairs */}
                    <div className="tab-switcher">
                        <button
                            className={`tab-btn ${folder === 'meta' ? 'active' : ''}`}
                            onClick={() => setFolder('meta')}
                        >
                            ⚡ Individual Equipment
                        </button>
                        <button
                            className={`tab-btn ${folder === 'metapairs' ? 'active' : ''}`}
                            onClick={() => setFolder('metapairs')}
                        >
                            🔗 Equipment Pairs
                        </button>
                    </div>

                    {/* Season Selector */}
                    <label className="control-select">
                        Season:
                        <select
                            value={month}
                            onChange={(e) => setMonth(e.target.value)}
                            className="select-input"
                        >
                            {months.map((m) => (
                                <option key={m} value={m}>
                                    {m}
                                </option>
                            ))}
                        </select>
                    </label>

                    {/* Dataset Size Selector */}
                    {files.length > 1 && (
                        <label className="control-select">
                            Sample:
                            <select
                                value={fileName}
                                onChange={(e) => setFileName(e.target.value)}
                                className="select-input"
                            >
                                {files.map((f) => (
                                    <option key={f} value={f}>
                                        {formatFileName(f)}
                                    </option>
                                ))}
                            </select>
                        </label>
                    )}
                </div>

                {/* Search */}
                <input
                    type="text"
                    placeholder="Search equipment..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="search-input"
                />
            </div>

            {/* Hero Quick Filter Chips */}
            {heroNames.length > 0 && (
                <div className="hero-chips-bar">
                    <button
                        className={`chip-btn ${selectedHero === 'All' ? 'active' : ''}`}
                        onClick={() => setSelectedHero('All')}
                    >
                        All Heroes ({heroNames.length})
                    </button>
                    {heroNames.map((name) => {
                        const theme = HERO_THEMES[name] || DEFAULT_THEME
                        return (
                            <button
                                key={name}
                                className={`chip-btn ${selectedHero === name ? 'active' : ''}`}
                                onClick={() => setSelectedHero(name)}
                            >
                                {theme.icon} {name}
                            </button>
                        )
                    })}
                </div>
            )}

            {/* Error & Loading States */}
            {error && <div className="error-banner">⚠️ Failed to load meta data: {error}</div>}
            {!data && !error && (
                <div className="state-container">
                    <div className="spinner"></div>
                    <p>Loading Clash of Clans meta data...</p>
                </div>
            )}

            {/* Heroes Grid */}
            {data && (
                <div className="heroes-grid">
                    {filteredHeroes.map(([hero, items]) => (
                        <HeroCard
                            key={hero}
                            hero={hero}
                            items={items}
                            searchQuery={searchQuery}
                        />
                    ))}
                </div>
            )}
                </>
            )}
        </div>
    )
}
