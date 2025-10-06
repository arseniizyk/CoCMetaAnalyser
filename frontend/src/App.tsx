import { useEffect, useState, useRef } from 'react'
import './App.css'
import { Chart } from 'chart.js/auto'

const dataFiles = import.meta.glob('/public/*/*/*.json', { eager: true })

const structure: Record<string, Record<string, string[]>> = {}

for (const path of Object.keys(dataFiles)) {
    const parts = path.split('/')
    const folder = parts.at(-3)!
    const month = parts.at(-2)!
    const file = parts.at(-1)!
    if (!structure[folder]) structure[folder] = {}
    if (!structure[folder][month]) structure[folder][month] = []
    structure[folder][month].push(file)
}

type HeroData = Record<string, number | string>
type JsonRoot = Record<string, HeroData>

function HeroChart({ hero, items }: { hero: string; items: HeroData }) {
    const ref = useRef<HTMLCanvasElement | null>(null)

    useEffect(() => {
        if (!ref.current) return

        const chart = new Chart(ref.current, {
            type: 'bar',
            data: {
                labels: Object.keys(items),
                datasets: [
                    { label: hero, data: Object.values(items).map(Number) },
                ],
            },
            options: {
                indexAxis: 'y',
                plugins: { legend: { display: false } },
                responsive: true,
                maintainAspectRatio: false,
            },
        })

        return () => chart.destroy()
    }, [hero, items])

    return (
        <div className="hero-card">
            <h3>{hero}</h3>
            <div style={{ height: 260 }}>
                <canvas ref={ref}></canvas>
            </div>
        </div>
    )
}

export default function App() {
    const [folder, setFolder] = useState(Object.keys(structure)[0])
    const [month, setMonth] = useState(Object.keys(structure[folder])[0])
    const [fileName, setFileName] = useState(structure[folder][month][0])
    const [data, setData] = useState<JsonRoot | null>(null)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        const path = `${import.meta.env.BASE_URL}${folder}/${month}/${fileName}`
        setError(null)
        fetch(path)
            .then((r) => {
                if (!r.ok) throw new Error(`Fetch ${path} failed: ${r.status}`)
                return r.json()
            })
            .then(setData)
            .catch((err) => setError(String(err)))
    }, [folder, month, fileName])

    return (
        <div style={{ padding: 16 }}>
            <div className="controls">
                <label>
                    Source:
                    <select
                        value={folder}
                        onChange={(e) => {
                            const newFolder = e.target.value
                            setFolder(newFolder)
                            const firstMonth = Object.keys(
                                structure[newFolder]
                            )[0]
                            setMonth(firstMonth)
                            setFileName(structure[newFolder][firstMonth][0])
                        }}
                    >
                        {Object.keys(structure).map((f) => (
                            <option key={f}>{f}</option>
                        ))}
                    </select>
                </label>

                <label>
                    Month:
                    <select
                        value={month}
                        onChange={(e) => {
                            const newMonth = e.target.value
                            setMonth(newMonth)
                            setFileName(structure[folder][newMonth][0])
                        }}
                    >
                        {Object.keys(structure[folder]).map((m) => (
                            <option key={m}>{m}</option>
                        ))}
                    </select>
                </label>

                <label>
                    File:
                    <select
                        value={fileName}
                        onChange={(e) => setFileName(e.target.value)}
                    >
                        {structure[folder][month].map((f) => (
                            <option key={f}>{f}</option>
                        ))}
                    </select>
                </label>
            </div>

            {error && <div style={{ color: 'crimson' }}>Error: {error}</div>}
            {!data && !error && <div>Loading...</div>}

            {data && (
                <div className="grid">
                    {Object.entries(data).map(([hero, items]) => (
                        <HeroChart key={hero} hero={hero} items={items} />
                    ))}
                </div>
            )}
        </div>
    )
}
