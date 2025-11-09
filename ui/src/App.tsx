import { FormEvent, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  createServer,
  deleteServer,
  listServers,
  ServerInfo,
  ServerMetadata,
  ServerType,
  startServer,
  stopServer,
} from './api';
import './App.css';

interface CreateFormState {
  worldName: string;
  type: ServerType;
  version: string;
  maxMemory: string;
  port: string;
  mods: string;
}

const initialCreateState: CreateFormState = {
  worldName: '',
  type: 'vanilla',
  version: '',
  maxMemory: '3G',
  port: '25565',
  mods: '',
};

export default function App() {
  const [createForm, setCreateForm] = useState<CreateFormState>(initialCreateState);
  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: () =>
      createServer({
        worldName: createForm.worldName.trim(),
        type: createForm.type,
        version: createForm.version.trim() || undefined,
        maxMemory: createForm.maxMemory.trim() || undefined,
        port: createForm.port.trim() || undefined,
        curseForgeMods: createForm.mods
          .split(',')
          .map((m) => m.trim())
          .filter(Boolean),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['servers'] });
    },
  });

  const startMutation = useMutation({
    mutationFn: (worldName: string) => startServer(worldName.trim()),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['servers'] });
    },
  });

  const stopMutation = useMutation({
    mutationFn: (worldName: string) => stopServer(worldName.trim()),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['servers'] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (worldName: string) => deleteServer(worldName.trim()),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['servers'] });
    },
  });

  const serversQuery = useQuery({
    queryKey: ['servers'],
    queryFn: listServers,
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

  return (
    <main className="app-shell">
      <header>
        <h1>mcrun UI</h1>
        <p>Provision, stop, delete, and inspect Minecraft servers through the HTTP API.</p>
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
            Port
            <input
              type="text"
              value={createForm.port}
              onChange={(e) => setCreateForm({ ...createForm, port: e.target.value })}
              placeholder="25565"
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
        <div className="servers-header">
          <h2>Existing servers</h2>
          <button
            type="button"
            disabled={serversQuery.isRefetching}
            onClick={() => serversQuery.refetch()}
          >
            {serversQuery.isRefetching ? 'Refreshing…' : 'Refresh'}
          </button>
        </div>
        <div className="card">
          {serversQuery.isLoading && <p className="info">Loading servers…</p>}
          {serversQuery.isError && (
            <p className="error">{(serversQuery.error as Error).message}</p>
          )}
          {serversQuery.isSuccess && serversQuery.data.length === 0 && (
            <p className="info">No servers found.</p>
          )}
          {serversQuery.isSuccess && serversQuery.data.length > 0 && (
            <div className="server-table">
              <div className="server-table__header">
                <span>World</span>
                <span>Type</span>
                <span>Version</span>
                <span>Port</span>
                <span>Mods</span>
                <span>Status</span>
                <span>Compose</span>
                <span>Actions</span>
              </div>
              {serversQuery.data.map((server) => (
                <ServerRow
                  key={server.worldName}
                  server={server}
                  onStart={() => startMutation.mutate(server.worldName)}
                  onStop={() => stopMutation.mutate(server.worldName)}
                  onDelete={() => deleteMutation.mutate(server.worldName)}
                  startBusy={
                    startMutation.isPending && startMutation.variables === server.worldName
                  }
                  stopBusy={stopMutation.isPending && stopMutation.variables === server.worldName}
                  deleteBusy={
                    deleteMutation.isPending && deleteMutation.variables === server.worldName
                  }
                />
              ))}
            </div>
          )}
        </div>
      </section>
    </main>
  );
}

interface ServerRowProps {
  server: ServerInfo;
  onStart: () => void;
  onStop: () => void;
  onDelete: () => void;
  startBusy: boolean;
  stopBusy: boolean;
  deleteBusy: boolean;
}

function ServerRow({
  server,
  onStart,
  onStop,
  onDelete,
  startBusy,
  stopBusy,
  deleteBusy,
}: ServerRowProps) {
  const meta: ServerMetadata | undefined = server.metadata;
  const mods =
    meta?.mods && meta.mods.length > 0
      ? meta.mods.map((m) => m.name || `#${m.curseforge?.projectId ?? '-'}`).join(', ')
      : '—';

  return (
    <div className="server-table__row">
      <span>
        <strong>{server.worldName}</strong>
        <small>{server.path}</small>
      </span>
      <span>{meta?.type ?? '—'}</span>
      <span>{meta?.version ?? '—'}</span>
      <span>{meta?.port ?? '25565'}</span>
      <span className="mods">{mods}</span>
      <span>{server.status ?? 'unknown'}</span>
      <span>{server.hasCompose ? 'yes' : 'no'}</span>
      <span className="action-buttons">
        <button
          type="button"
          className="icon-button play"
          onClick={onStart}
          disabled={startBusy}
          title="Start server"
        >
          {startBusy ? SpinnerIcon : PlayIcon}
        </button>
        <button
          type="button"
          className="icon-button"
          onClick={onStop}
          disabled={stopBusy}
          title="Stop server"
        >
          {stopBusy ? SpinnerIcon : StopIcon}
        </button>
        <button
          type="button"
          className="icon-button danger"
          onClick={onDelete}
          disabled={deleteBusy}
          title="Delete server"
        >
          {deleteBusy ? SpinnerIcon : TrashIcon}
        </button>
      </span>
    </div>
  );
}

const PlayIcon = (
  <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
    <path fill="currentColor" d="M8 5v14l11-7z" />
  </svg>
);

const StopIcon = (
  <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
    <path fill="currentColor" d="M6 6h12v12H6z" />
  </svg>
);

const TrashIcon = (
  <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
    <path
      fill="currentColor"
      d="M9 3h6l1 1h4v2H4V4h4l1-1zm1 6h2v8h-2V9zm4 0h2v8h-2V9zM6 7h12v12a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2V7z"
    />
  </svg>
);

const SpinnerIcon = (
  <svg viewBox="0 0 24 24" width="18" height="18" className="spin" aria-hidden="true">
    <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" opacity="0.25" />
    <path
      fill="currentColor"
      d="M22 12a10 10 0 0 1-10 10v-4a6 6 0 0 0 6-6h4z"
    />
  </svg>
);
