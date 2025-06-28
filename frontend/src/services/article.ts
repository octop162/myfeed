import api from './api';

export interface Article {
  id: string;
  feed_id: string;
  title: string;
  content?: string;
  url: string;
  published_at?: string;
  is_read: boolean;
  is_later: boolean;
  created_at: string;
}

export const articleService = {
  getAllArticles: async (): Promise<Article[]> => {
    const response = await api.get<Article[]>('/articles');
    return response.data;
  },

  getArticleById: async (id: string): Promise<Article> => {
    const response = await api.get<Article>(`/articles/${id}`);
    return response.data;
  },

  updateArticleStatus: async (id: string, is_read: boolean, is_later: boolean): Promise<Article> => {
    const response = await api.put<Article>(`/articles/${id}/status`, { is_read, is_later });
    return response.data;
  },

  getLaterArticles: async (): Promise<Article[]> => {
    const response = await api.get<Article[]>('/articles/later');
    return response.data;
  },
};
