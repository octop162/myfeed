import api from './api';

export interface Feed {
  id: string;
  name: string;
  url: string;
  plugin_type: string;
  folder_id?: string;
  update_interval: number;
  last_updated?: string;
  created_at: string;
}

export const feedService = {
  getAllFeeds: async (): Promise<Feed[]> => {
    const response = await api.get<Feed[]>('/feeds');
    return response.data;
  },

  getFeedById: async (id: string): Promise<Feed> => {
    const response = await api.get<Feed>(`/feeds/${id}`);
    return response.data;
  },

  createFeed: async (feed: Omit<Feed, 'id' | 'created_at' | 'last_updated'>): Promise<Feed> => {
    const response = await api.post<Feed>('/feeds', feed);
    return response.data;
  },

  updateFeed: async (id: string, feed: Omit<Feed, 'id' | 'created_at' | 'last_updated'>): Promise<Feed> => {
    const response = await api.put<Feed>(`/feeds/${id}`, feed);
    return response.data;
  },

  deleteFeed: async (id: string): Promise<void> => {
    await api.delete(`/feeds/${id}`);
  },
};
