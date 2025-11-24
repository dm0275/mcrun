import { Dispatch, FormEvent, SetStateAction, useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  createServer,
  deleteServer,
  listServers,
  ModSpec,
  ServerInfo,
  ServerMetadata,
  ServerType,
  startServer,
  stopServer,
  updateServer,
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

interface UpdateServerInput {
  worldName: string;
  maxMemory?: string;
  mods: string[];
}

interface EditFormState {
  maxMemory: string;
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
  const [portError, setPortError] = useState<string | null>(null);
  const [manifestError, setManifestError] = useState<string | null>(null);

  const createMutation = useMutation({
    mutationFn: () =>
      createServer({
        worldName: createForm.worldName.trim(),
        type: createForm.type,
        version: createForm.version.trim() || undefined,
        maxMemory: createForm.maxMemory.trim() || undefined,
        port: createForm.port.trim() || undefined,
        curseForgeMods: parseModsInput(createForm.mods),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['servers'] });
      setPortError(null);
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

  const updateMutation = useMutation({
    mutationFn: (payload: UpdateServerInput) =>
      updateServer({
        worldName: payload.worldName,
        maxMemory: payload.maxMemory,
        curseForgeMods: payload.mods,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['servers'] });
    },
  });

  const serversQuery = useQuery({
    queryKey: ['servers'],
    queryFn: listServers,
  });

  const findPortConflict = useMemo(() => {
    return (port: string, excludeWorld?: string) => {
      if (!serversQuery.data) {
        return undefined;
      }
      const normalized = port.trim();
      return serversQuery.data.find(
        (server) =>
          server.status === 'running' &&
          server.metadata?.port === normalized &&
          server.worldName !== excludeWorld,
      );
    };
  }, [serversQuery.data]);

  const onSubmitCreate = (event: FormEvent) => {
    event.preventDefault();
    setPortError(null);
    if (!createForm.worldName.trim()) {
      return;
    }
    const desiredPort = createForm.port.trim() || '25565';
    const conflict = findPortConflict(desiredPort);
    if (conflict) {
      setPortError(`Port ${desiredPort} is already used by ${conflict.worldName}. Stop it first.`);
      return;
    }
    createMutation.mutate(undefined, {
      onSuccess: () => setCreateForm(initialCreateState),
    });
  };

  return (
    <main className="app-shell">
      <header className="app-header">
        <div className="header-content">
          <h1 className="app-title">
            <span className="title-icon">⚙️</span>
            MCRUN
          </h1>
          <p className="header-description">Manage your Minecraft servers with ease</p>
        </div>
      </header>

      <section>
        <h2 className="section-title">
          <span className="section-icon">➕</span>
          Create New Server
        </h2>
        <form className="card create-form" onSubmit={onSubmitCreate}>
          <label>
            Import CurseForge manifest (.json)
            <input
              type="file"
              accept="application/json,.json"
              onChange={(e) => {
                const file = e.target.files?.[0];
                if (!file) {
                  return;
                }
                parseManifestFile(file, setCreateForm, setManifestError);
              }}
            />
            <small className="helper-text">
              Prefills world, server type, and mods using a CurseForge modpack manifest.
            </small>
          </label>
          {manifestError && <p className="error">{manifestError}</p>}

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
              onChange={(e) => {
                setCreateForm({ ...createForm, port: e.target.value });
                if (portError) {
                  setPortError(null);
                }
              }}
              placeholder="25565"
            />
          </label>
          {portError && <p className="error">{portError}</p>}

          <label>
            CurseForge mods (comma separated `projectId[:fileId]@version`)
            <textarea
              rows={3}
              value={createForm.mods}
              onChange={(e) => setCreateForm({ ...createForm, mods: e.target.value })}
              placeholder="238222:6570130@1.20.1, 306612"
            />
          </label>

          <button type="submit" className="btn-primary" disabled={createMutation.isPending}>
            {createMutation.isPending ? (
              <>
                <span className="spinner-small"></span>
                Creating server…
              </>
            ) : (
              <>
                <span>🚀</span>
                Create server
              </>
            )}
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
          <h2 className="section-title">
            <span className="section-icon">🖥️</span>
            Server Dashboard
          </h2>
          <button
            type="button"
            className="btn-refresh"
            disabled={serversQuery.isRefetching}
            onClick={() => serversQuery.refetch()}
            title="Refresh server list"
          >
            {serversQuery.isRefetching ? (
              <>
                <span className="spinner-small"></span>
                Refreshing…
              </>
            ) : (
              <>
                <span>🔄</span>
                Refresh
              </>
            )}
          </button>
        </div>
        <div className="card">
          {serversQuery.isLoading && (
            <div className="empty-state">
              <span className="empty-icon">⏳</span>
              <p className="info">Loading servers…</p>
            </div>
          )}
          {serversQuery.isError && (
            <div className="empty-state">
              <span className="empty-icon">❌</span>
              <p className="error">{(serversQuery.error as Error).message}</p>
            </div>
          )}
          {serversQuery.isSuccess && serversQuery.data.length === 0 && (
            <div className="empty-state">
              <span className="empty-icon">📦</span>
              <p className="info">No servers found. Create your first server to get started!</p>
            </div>
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
                  portConflict={
                    server.metadata?.port
                      ? findPortConflict(server.metadata.port, server.worldName)
                      : undefined
                  }
                  onUpdate={(payload, options) =>
                    updateMutation.mutate(payload, {
                      onSuccess: () => {
                        options?.onSuccess?.();
                      },
                    })
                  }
                  updateBusy={
                    updateMutation.isPending &&
                    updateMutation.variables?.worldName === server.worldName
                  }
                  updateError={
                    updateMutation.isError &&
                    updateMutation.variables?.worldName === server.worldName
                      ? updateMutation.error instanceof Error
                        ? updateMutation.error.message
                        : 'Failed to update server'
                      : null
                  }
                  resetUpdate={updateMutation.reset}
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
   onUpdate: (payload: UpdateServerInput, options?: { onSuccess?: () => void }) => void;
  startBusy: boolean;
  stopBusy: boolean;
  deleteBusy: boolean;
  updateBusy: boolean;
  updateError: string | null;
  resetUpdate: () => void;
  portConflict?: ServerInfo;
}

function ServerRow({
  server,
  onStart,
  onStop,
  onDelete,
  onUpdate,
  startBusy,
  stopBusy,
  deleteBusy,
  updateBusy,
  updateError,
  resetUpdate,
  portConflict,
}: ServerRowProps) {
  const meta: ServerMetadata | undefined = server.metadata;
  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState<EditFormState>(() => buildEditFormState(meta));
  const modsDisplay = formatMods(meta?.mods);

  const statusClass = server.status === 'running' ? 'status-running' : server.status === 'stopped' ? 'status-stopped' : 'status-unknown';
  const typeClass = meta?.type ? `type-badge type-${meta.type}` : 'type-badge';

  const startDisabled = startBusy || Boolean(portConflict);
  const startTitle = portConflict
    ? `Port ${meta?.port ?? '25565'} used by ${portConflict.worldName}`
    : 'Start server';

  const handleToggleEdit = () => {
    resetUpdate();
    if (isEditing) {
      setIsEditing(false);
      setEditForm(buildEditFormState(meta));
      return;
    }
    setEditForm(buildEditFormState(meta));
    setIsEditing(true);
  };

  const handleUpdateSubmit = (event: FormEvent) => {
    event.preventDefault();
    const normalizedMax = editForm.maxMemory.trim();
    const modsList = parseModsInput(editForm.mods);

    onUpdate(
      {
        worldName: server.worldName,
        maxMemory: normalizedMax || undefined,
        mods: modsList,
      },
      {
        onSuccess: () => {
          setIsEditing(false);
          setEditForm({
            maxMemory: normalizedMax,
            mods: editForm.mods.trim(),
          });
        },
      },
    );
  };

  return (
    <>
      <div className="server-table__row">
        <span>
          <strong className="world-name">{server.worldName}</strong>
        </span>
        <span><span className={typeClass}>{meta?.type ?? '—'}</span></span>
        <span className="version">{meta?.version ?? '—'}</span>
        <span className="port">{meta?.port ?? '25565'}</span>
        <span className="mods" title={modsDisplay.title || undefined}>{modsDisplay.display}</span>
        <span><span className={`status-badge ${statusClass}`}>{server.status ?? 'unknown'}</span></span>
        <span>{server.hasCompose ? <span className="compose-badge">✓</span> : '—'}</span>
        <span className="action-buttons">
          <button
            type="button"
            className="icon-button play"
            onClick={onStart}
            disabled={startDisabled}
            title={startTitle}
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
            className="icon-button"
            onClick={handleToggleEdit}
            disabled={updateBusy}
            title="Edit server config"
          >
            {updateBusy ? SpinnerIcon : EditIcon}
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
      {isEditing && (
        <div className="server-edit-panel">
          <form className="server-edit-form" onSubmit={handleUpdateSubmit}>
            <div className="server-edit-grid">
              <label>
                Max memory
                <input
                  type="text"
                  value={editForm.maxMemory}
                  onChange={(e) => setEditForm({ ...editForm, maxMemory: e.target.value })}
                  placeholder="3G"
                />
              </label>
              <label>
                CurseForge mods (comma separated `projectId[:fileId]@version`)
                <textarea
                  rows={3}
                  value={editForm.mods}
                  onChange={(e) => setEditForm({ ...editForm, mods: e.target.value })}
                  placeholder="238222:6570130@1.20.1, 306612"
                />
              </label>
            </div>
            {updateError && <p className="error">{updateError}</p>}
            <div className="server-edit-actions">
              <button type="submit" className="btn-primary" disabled={updateBusy}>
                {updateBusy ? (
                  <>
                    <span className="spinner-small"></span>
                    Saving…
                  </>
                ) : (
                  <>
                    <span>💾</span>
                    Save changes
                  </>
                )}
              </button>
              <button
                type="button"
                className="btn-secondary"
                onClick={handleToggleEdit}
                disabled={updateBusy}
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}
    </>
  );
}

function parseManifestFile(
  file: File,
  updateForm: Dispatch<SetStateAction<CreateFormState>>,
  setError: (message: string | null) => void,
) {
  setError(null);
  file
    .text()
    .then((content) => {
      let manifest: any;
      try {
        manifest = JSON.parse(content);
      } catch (err) {
        throw new Error('Invalid JSON manifest file');
      }

      const name = typeof manifest?.name === 'string' ? manifest.name : '';
      const minecraftVersion =
        typeof manifest?.minecraft?.version === 'string' ? manifest.minecraft.version : '';
      const loaderId =
        manifest?.minecraft?.modLoaders?.find((l: any) => l?.primary)?.id ??
        manifest?.minecraft?.modLoaders?.[0]?.id ??
        '';

      const inferredType: ServerType = inferServerType(loaderId);
      const normalizedWorld = slugifyWorldName(name);
      const manifestMods = Array.isArray(manifest?.files)
        ? manifest.files
            .map((entry: any) => {
              const projectId = entry?.projectID ?? entry?.projectId;
              const fileId = entry?.fileID ?? entry?.fileId;
              if (!projectId) {
                return null;
              }
              const versionPart = minecraftVersion ? `@${minecraftVersion}` : '';
              const filePart = fileId ? `:${fileId}` : '';
              return `${projectId}${filePart}${versionPart}`;
            })
            .filter(Boolean)
        : [];

      updateForm((prev) => ({
        ...prev,
        worldName: normalizedWorld || prev.worldName,
        type: inferredType ?? prev.type,
        version:
          minecraftVersion && inferredType !== 'vanilla'
            ? `${inferredType}-${minecraftVersion}`
            : prev.version || minecraftVersion,
        mods: manifestMods.length > 0 ? Array.from(new Set(manifestMods)).join(', ') : prev.mods,
      }));
    })
    .catch((err: unknown) => {
      setError(err instanceof Error ? err.message : 'Failed to parse manifest file');
    });
}

function inferServerType(loaderId: string): ServerType {
  const normalized = loaderId.toLowerCase();
  if (normalized.includes('forge')) {
    return 'forge';
  }
  if (normalized.includes('fabric')) {
    return 'fabric';
  }
  return 'vanilla';
}

function slugifyWorldName(name: string): string {
  const trimmed = name.trim().toLowerCase();
  if (!trimmed) {
    return '';
  }
  const slug = trimmed.replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)+/g, '');
  return slug || trimmed;
}

const MOD_DISPLAY_LIMIT = 6;
function formatMods(mods?: ModSpec[]) {
  if (!mods || mods.length === 0) {
    return { display: '—', title: '' };
  }

  const names = mods.map((m) => m.name || (m.curseforge?.projectId ? `#${m.curseforge.projectId}` : 'mod'));
  const title = names.join(', ');
  const visible = names.slice(0, MOD_DISPLAY_LIMIT);
  const remaining = names.length - visible.length;
  const display = remaining > 0 ? `${visible.join(', ')} (+${remaining} more)` : visible.join(', ');

  return { display, title };
}

function parseModsInput(modsText: string) {
  return modsText
    .split(',')
    .map((m) => m.trim())
    .filter(Boolean);
}

function modsToInput(mods?: ModSpec[]) {
  if (!mods || mods.length === 0) {
    return '';
  }

  return mods
    .map((mod) => {
      const curse = mod.curseforge;
      if (curse?.projectId) {
        const filePart = curse.fileId ? `:${curse.fileId}` : '';
        const versionPart = curse.gameVersion ? `@${curse.gameVersion}` : '';
        return `${curse.projectId}${filePart}${versionPart}`;
      }
      return '';
    })
    .filter(Boolean)
    .join(', ');
}

function buildEditFormState(meta?: ServerMetadata): EditFormState {
  return {
    maxMemory: (meta?.maxMemory || meta?.minMemory || '').trim(),
    mods: modsToInput(meta?.mods),
  };
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

const EditIcon = (
  <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
    <path
      fill="currentColor"
      d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zm2.92 2.17H5v-0.92l9.06-9.06 0.92 0.92L5.92 19.42zM20.71 5.63l-2.34-2.34a1.003 1.003 0 0 0-1.42 0l-1.83 1.83 3.75 3.75 1.84-1.84a1.003 1.003 0 0 0 0-1.4z"
    />
  </svg>
);
