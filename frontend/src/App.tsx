import { useEffect, useState } from 'react'
import './App.css'
import { Chart } from 'chart.js/auto'

type HeroData = Record<string, number | string>
type JsonRoot = Record<string, HeroData>

const MONTHS: Record<string, string[]> = {
    meta: [
        '2025-01',
        '2025-02',
        '2025-03',
        '2025-04',
        '2025-05',
        '2025-06',
        '2025-07',
        '2025-08'
    ],
    metapairs: [
        '2025-02',
        '2025-03',
        '2025-04',
        '2025-05',
        '2025-06',
        '2025-07',
        '2025-08'
    ],
}

const FILES: Record<string, string[]> = {
    meta: ['meta10k.json'],
    metapairs: ['meta10k.json', 'meta100.json'],
}

function HeroChart({ hero, items }: { hero: string; items: HeroData }) {
    const canvasId = `chart-${hero.replace(/\s+/g, '_')}`

    useEffect(() => {
        const canvas = document.getElementById(
            canvasId
        ) as HTMLCanvasElement | null
        if (!canvas) return

        // destroy old chart
        // @ts-ignore
        if (canvas.__chart) {
            // @ts-ignore
            canvas.__chart.destroy()
            // @ts-ignore
            canvas.__chart = null
        }

        const labels = Object.keys(items)
        const data = Object.values(items).map((v) => Number(v || 0))

        // @ts-ignore
        canvas.__chart = new Chart(canvas, {
            type: 'bar',
            data: { labels, datasets: [{ label: hero, data, borderWidth: 1 }] },
            options: {
                indexAxis: 'y',
                responsive: true,
                maintainAspectRatio: false,
                plugins: { legend: { display: false } },
                scales: {
                    x: {
                        beginAtZero: true,
                        ticks: {
                            autoSkip: false,
                            maxRotation: 90,
                            minRotation: 30,
                        },
                    },
                    y: { ticks: { autoSkip: false } },
                },
            },
        })

        return () => {
            // @ts-ignore
            if (canvas.__chart) {
                // @ts-ignore
                canvas.__chart.destroy()
                // @ts-ignore
                canvas.__chart = null
            }
        }
    }, [canvasId, JSON.stringify(items)])

    return (
        <div className="hero-card">
            <h3>{hero}</h3>
            <div style={{ height: 260 }}>
                <canvas id={canvasId}></canvas>
            </div>
        </div>
    )
}

export default function App() {
    const [folder, setFolder] = useState<'meta' | 'metapairs'>('meta')
    const [month, setMonth] = useState(MONTHS['meta'][0])
    const [fileName, setFileName] = useState(FILES['meta'][0])
    const [data, setData] = useState<JsonRoot | null>(null)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        const months = MONTHS[folder]
        const files = FILES[folder]

        const effectiveMonth = months.includes(month) ? month : months[0]
        const effectiveFile = files.includes(fileName) ? fileName : files[0]

        if (effectiveMonth !== month) setMonth(effectiveMonth)
        if (effectiveFile !== fileName) setFileName(effectiveFile)

        setError(null)
        setData(null)
        
        const path = `${
            import.meta.env.BASE_URL
        }${folder}/${effectiveMonth}/${effectiveFile}`

        fetch(path)
            .then((r) => {
                if (!r.ok) throw new Error(`Fetch ${path} failed: ${r.status}`)
                return r.json()
            })
            .then((j: JsonRoot) => setData(j))
            .catch((err) => {
                console.warn(err)
                setError(String(err))
            })
    }, [folder, month, fileName])

    return (
        <div style={{ padding: 16 }}>
            <div className="controls">
                <label>
                    Source:
                    <select
                        value={folder}
                        onChange={(e) =>
                            setFolder(e.target.value as 'meta' | 'metapairs')
                        }
                    >
                        <option value="meta">meta</option>
                        <option value="metapairs">metapairs</option>
                    </select>
                </label>

                <label>
                    Month:
                    <select
                        value={month}
                        onChange={(e) => setMonth(e.target.value)}
                    >
                        {MONTHS[folder].map((m) => (
                            <option key={m} value={m}>
                                {m}
                            </option>
                        ))}
                    </select>
                </label>

                <label>
                    File:
                    <select
                        value={fileName}
                        onChange={(e) => setFileName(e.target.value)}
                    >
                        {FILES[folder].map((f) => (
                            <option key={f} value={f}>
                                {f}
                            </option>
                        ))}
                    </select>
                </label>
            </div>

            {error && (
                <div style={{ color: 'crimson', marginTop: 12 }}>
                    Error: {error}
                </div>
            )}
            {!data && !error && <div style={{ marginTop: 12 }}>Loading...</div>}

            {data && (
                <div className="grid">
                    {Object.keys(data).map((hero) => (
                        <HeroChart key={hero} hero={hero} items={data[hero]} />
                    ))}
                </div>
            )}
        </div>
    )
}
