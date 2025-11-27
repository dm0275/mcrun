const API_BASE = '/api';

export type ServerType = 'vanilla' | 'forge' | 'fabric';

export interface CreateServerPayload {
  worldName: string;
  type: ServerType;
  version?: string;
  maxMemory?: string;
  minMemory?: string;
  port?: string;
  seed?: string;
  curseForgeMods?: string[];
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers || {}),
    },
    ...init,
  });

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`;
    try {
      const data = await response.json();
      message = data.error ?? message;
    } catch (_) {
      // ignore JSON parse errors
    }
    throw new Error(message);
  }

  if (response.status === 204) {
    return {} as T;
  }

  return (await response.json()) as T;
}

export function createServer(payload: CreateServerPayload) {
  const body = {
    worldName: payload.worldName,
    type: payload.type,
    version: payload.version,
    maxMemory: payload.maxMemory,
    minMemory: payload.minMemory,
    port: payload.port,
    seed: payload.seed,
    mods: payload.curseForgeMods?.map((specifier) => ({
      source: 'curseforge',
      curseforge: parseModSpecifier(specifier),
    })),
  };
  return request<{ status: string; worldName: string }>(`/servers`, {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export interface ServerInfo {
  worldName: string;
  path: string;
  hasCompose: boolean;
  status?: string;
  metadata?: ServerMetadata;
}

export interface ServerMetadata {
  worldName: string;
  type: string;
  version: string;
  port: string;
  image: string;
  maxMemory?: string;
  minMemory?: string;
  mods?: ModSpec[];
  enableRcon?: boolean;
  rconPort?: string;
  rconPassword?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ModSpec {
  name?: string;
  source: string;
  url?: string;
  downloadUrl?: string;
  hash?: string;
  curseforge?: {
    projectId: number;
    fileId?: number;
    gameVersion?: string;
    loader?: string;
  };
}

export function listServers() {
  return request<ServerInfo[]>(`/servers`);
}

export function stopServer(worldName: string) {
  return request<{ status: string }>(`/servers/${encodeURIComponent(worldName)}/stop`, {
    method: 'POST',
  });
}

export function startServer(worldName: string) {
  return request<{ status: string }>(`/servers/${encodeURIComponent(worldName)}/start`, {
    method: 'POST',
  });
}

export function deleteServer(worldName: string) {
  return request<{ status: string }>(`/servers/${encodeURIComponent(worldName)}`, {
    method: 'DELETE',
  });
}

export interface UpdateServerPayload {
  worldName: string;
  maxMemory?: string;
  minMemory?: string;
  curseForgeMods?: string[];
}

export function updateServer(payload: UpdateServerPayload) {
  const body: any = {};
  if (payload.maxMemory !== undefined) {
    body.maxMemory = payload.maxMemory;
  }
  if (payload.minMemory !== undefined) {
    body.minMemory = payload.minMemory;
  }
  if (payload.curseForgeMods) {
    body.mods = payload.curseForgeMods.map((specifier) => ({
      source: 'curseforge',
      curseforge: parseModSpecifier(specifier),
    }));
  }

  return request<{ status: string; worldName: string }>(
    `/servers/${encodeURIComponent(payload.worldName)}`,
    {
      method: 'PATCH',
      body: JSON.stringify(body),
    },
  );
}

export function sendRconCommand(worldName: string, command: string) {
  return request<{ response: string }>(`/servers/${encodeURIComponent(worldName)}/rcon`, {
    method: 'POST',
    body: JSON.stringify({ command }),
  });
}

function parseModSpecifier(spec: string) {
  const [projectAndFile, version] = spec.split('@');
  const [projectIdStr, fileIdStr] = projectAndFile.split(':');
  return {
    projectId: Number(projectIdStr.trim()),
    fileId: fileIdStr ? Number(fileIdStr.trim()) : undefined,
    gameVersion: version?.trim() ?? '',
  };
}
