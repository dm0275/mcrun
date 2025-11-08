import { FormEvent, useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { createServer, deleteServer, stopServer, ServerType } from './api';
import './App.css';

interface CreateFormState {
  worldName: string;
  type: ServerType;
  version: string;
  maxMemory: string;
  mods: string;
}

const initialCreateState: CreateFormState = {
  worldName: '',
  type: 'vanilla',
  version: '',
  maxMemory: '3G',
  mods: '',
};

export default function App() {
  const [createForm, setCreateForm] = useState<CreateFormState>(initialCreateState);
  const [manageWorldName, setManageWorldName] = useState('');

  const createMutation = useMutation({
    mutationFn: () =>
      createServer({
        worldName: createForm.worldName.trim(),
        type: createForm.type,
        version: createForm.version.trim() || undefined,
        maxMemory: createForm.maxMemory.trim() || undefined,
        curseForgeMods: createForm.mods
          .split(',')
          .map((m) => m.trim())
          .filter(Boolean),
      }),
  });

  const stopMutation = useMutation({
    mutationFn: () => stopServer(manageWorldName.trim()),
  });

  const deleteMutation = useMutation({
    mutationFn: () => deleteServer(manageWorldName.trim()),
  });

  const onSubmitCreate = (event: FormEvent) => {
    event.preventDefault();
    if (!createForm.worldName.trim()) {
      return;
    }
    createMutation.mutate(undefined, {
      onSuccess: () => setCreateForm(initialCreateState),
    });
  };

  const busy =
    createMutation.isPending || stopMutation.isPending || deleteMutation.isPending;

  return (
    <main className="app-shell">
      <header>
        <h1>mcrun UI</h1>
        <p>Provision, stop, or delete Minecraft servers through the HTTP API.</p>
      </header>

      <section>
        <h2>Create server</h2>
        <form className="card" onSubmit={onSubmitCreate}>
          <label>
            World name
            <input
              type="text"
              value={createForm.worldName}
              onChange={(e) => setCreateForm({ ...createForm, worldName: e.target.value })}
              required
              placeholder="my-awesome-world"
            />
          </label>

          <label>
            Server type
            <select
              value={createForm.type}
              onChange={(e) =>
                setCreateForm({ ...createForm, type: e.target.value as ServerType })
              }
            >
              <option value="vanilla">Vanilla</option>
              <option value="forge">Forge</option>
              <option value="fabric">Fabric</option>
            </select>
          </label>

          <label>
            Version (optional)
            <input
              type="text"
              value={createForm.version}
              onChange={(e) => setCreateForm({ ...createForm, version: e.target.value })}
              placeholder="forge-1.20.1"
            />
          </label>

          <label>
            Max memory
            <input
              type="text"
              value={createForm.maxMemory}
              onChange={(e) => setCreateForm({ ...createForm, maxMemory: e.target.value })}
              placeholder="3G"
            />
          </label>

          <label>
            CurseForge mods (comma separated `projectId@version`)
            <textarea
              rows={3}
              value={createForm.mods}
              onChange={(e) => setCreateForm({ ...createForm, mods: e.target.value })}
              placeholder="238222@1.20.1, 306612"
            />
          </label>

          <button type="submit" disabled={createMutation.isPending}>
            {createMutation.isPending ? 'Creating…' : 'Create server'}
          </button>

          {createMutation.isError && (
            <p className="error">{createMutation.error.message}</p>
          )}
          {createMutation.isSuccess && (
            <p className="success">
              Server {createMutation.data.worldName} is {createMutation.data.status}
            </p>
          )}
        </form>
      </section>

      <section>
        <h2>Manage existing server</h2>
        <div className="card">
          <label>
            World name
            <input
              type="text"
              value={manageWorldName}
              onChange={(e) => setManageWorldName(e.target.value)}
              placeholder="my-awesome-world"
            />
          </label>
          <div className="actions">
            <button
              onClick={() => stopMutation.mutate()}
              disabled={!manageWorldName || stopMutation.isPending}
            >
              {stopMutation.isPending ? 'Stopping…' : 'Stop server'}
            </button>
            <button
              className="danger"
              onClick={() => deleteMutation.mutate()}
              disabled={!manageWorldName || deleteMutation.isPending}
            >
              {deleteMutation.isPending ? 'Deleting…' : 'Delete server'}
            </button>
          </div>

          {busy && <p className="info">Working…</p>}
          {stopMutation.isError && <p className="error">{stopMutation.error.message}</p>}
          {deleteMutation.isError && (
            <p className="error">{deleteMutation.error.message}</p>
          )}
          {stopMutation.isSuccess && (
            <p className="success">Server {manageWorldName} is stopping</p>
          )}
          {deleteMutation.isSuccess && (
            <p className="success">Server {manageWorldName} deleted</p>
          )}
        </div>
      </section>
    </main>
  );
}
