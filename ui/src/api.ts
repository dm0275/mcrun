const API_BASE = '/api';

export type ServerType = 'vanilla' | 'forge' | 'fabric';

export interface CreateServerPayload {
  worldName: string;
  type: ServerType;
  version?: string;
  maxMemory?: string;
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

export function stopServer(worldName: string) {
  return request<{ status: string }>(`/servers/${encodeURIComponent(worldName)}/stop`, {
    method: 'POST',
  });
}

export function deleteServer(worldName: string) {
  return request<{ status: string }>(`/servers/${encodeURIComponent(worldName)}`, {
    method: 'DELETE',
  });
}

function parseModSpecifier(spec: string) {
  const [projectId, version] = spec.split('@');
  return {
    projectId: Number(projectId.trim()),
    gameVersion: version?.trim() ?? '',
  };
}
