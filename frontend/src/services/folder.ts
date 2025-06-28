import api from './api';

export interface Folder {
  id: string;
  name: string;
  created_at: string;
}

export const folderService = {
  getAllFolders: async (): Promise<Folder[]> => {
    const response = await api.get<Folder[]>('/folders');
    return response.data;
  },

  getFolderById: async (id: string): Promise<Folder> => {
    const response = await api.get<Folder>(`/folders/${id}`);
    return response.data;
  },

  createFolder: async (name: string): Promise<Folder> => {
    const response = await api.post<Folder>('/folders', { name });
    return response.data;
  },

  updateFolder: async (id: string, name: string): Promise<Folder> => {
    const response = await api.put<Folder>(`/folders/${id}`, { name });
    return response.data;
  },

  deleteFolder: async (id: string): Promise<void> => {
    await api.delete(`/folders/${id}`);
  },
};
