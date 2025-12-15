import { useEffect, useMemo, useState, type FormEvent } from 'react'
import './App.css'

type Note = {
  id: string
  title: string
  content: string
  createdAt: string
  updatedAt: string
}

type CreateNoteRequest = {
  title: string
  content: string
}

async function api<T>(input: RequestInfo | URL, init?: RequestInit): Promise<T> {
  const res = await fetch(input, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {}),
    },
  })

  if (!res.ok) {
    const text = await res.text().catch(() => '')
    throw new Error(`${res.status} ${res.statusText}${text ? `: ${text}` : ''}`)
  }

  // 204 No Content
  if (res.status === 204) {
    return undefined as T
  }

  return (await res.json()) as T
}

function App() {
  const [notes, setNotes] = useState<Note[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')

  const sortedNotes = useMemo(() => {
    return [...notes].sort((a, b) => b.updatedAt.localeCompare(a.updatedAt))
  }, [notes])

  async function refresh() {
    setLoading(true)
    setError(null)
    try {
      const data = await api<Note[]>('/api/notes')
      setNotes(data)
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void refresh()
  }, [])

  async function createNote(e: FormEvent) {
    e.preventDefault()
    setError(null)

    const payload: CreateNoteRequest = {
      title: title.trim(),
      content,
    }

    try {
      const created = await api<Note>('/api/notes', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      setNotes((prev) => [created, ...prev])
      setTitle('')
      setContent('')
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
    }
  }

  return (
    <div style={{ maxWidth: 900, margin: '0 auto', padding: 24 }}>
      <h1>Notes</h1>

      <form onSubmit={createNote} style={{ display: 'grid', gap: 8, marginBottom: 24 }}>
        <input
          placeholder="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
        <textarea
          placeholder="Write a note..."
          value={content}
          onChange={(e) => setContent(e.target.value)}
          rows={6}
        />
        <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
          <button type="submit" disabled={!title.trim() && !content.trim()}>
            Create
          </button>
          <button type="button" onClick={() => void refresh()} disabled={loading}>
            Refresh
          </button>
          {loading ? <span>Loading…</span> : null}
        </div>
      </form>

      {error ? (
        <p style={{ color: 'crimson' }}>
          <strong>Error:</strong> {error}
        </p>
      ) : null}

      <ul style={{ listStyle: 'none', padding: 0, display: 'grid', gap: 12 }}>
        {sortedNotes.map((n) => (
          <li key={n.id} style={{ border: '1px solid #ddd', borderRadius: 8, padding: 12 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', gap: 12 }}>
              <strong>{n.title || '(untitled)'}</strong>
              <small>{new Date(n.updatedAt).toLocaleString()}</small>
            </div>
            {n.content ? <pre style={{ whiteSpace: 'pre-wrap', marginTop: 8 }}>{n.content}</pre> : null}
          </li>
        ))}
      </ul>

      {sortedNotes.length === 0 && !loading ? <p>No notes yet.</p> : null}
    </div>
  )
}

export default App
