import React from 'react';
import { Article, articleService } from '../services/article';
import { useSWRConfig } from 'swr';

interface ArticleListProps {
  articles: Article[];
}

const ArticleList: React.FC<ArticleListProps> = ({ articles }) => {
  const { mutate } = useSWRConfig();
  const handleToggleRead = async (article: Article) => {
    try {
      await articleService.updateArticleStatus(article.id, !article.is_read, article.is_later);
      mutate('/articles'); // 記事一覧を再フェッチ
      mutate('/articles/later'); // 後で見る記事一覧も再フェッチ
    } catch (error) {
      console.error('Failed to update article status:', error);
      alert('Failed to update article status.');
    }
  };

  const handleToggleLater = async (article: Article) => {
    try {
      await articleService.updateArticleStatus(article.id, article.is_read, !article.is_later);
      mutate('/articles'); // 記事一覧を再フェッチ
      mutate('/articles/later'); // 後で見る記事一覧も再フェッチ
    } catch (error) {
      console.error('Failed to update article status:', error);
      alert('Failed to update article status.');
    }
  };

  return (
    <ul className="space-y-4">
      {articles.map((article) => (
        <li key={article.id} className="p-4 border rounded-lg shadow-sm">
          <h2 className="text-xl font-semibold text-blue-600 hover:underline">
            <a href={article.url} target="_blank" rel="noopener noreferrer">
              {article.title}
            </a>
          </h2>
          {article.content && <p className="text-gray-700 mt-2">{article.content.substring(0, 150)}...</p>}
          <div className="text-sm text-gray-500 mt-2">
            <span>{new Date(article.published_at || article.created_at).toLocaleDateString()}</span>
            <span className="ml-4">{article.is_read ? 'Read' : 'Unread'}</span>
            <span className="ml-4">{article.is_later ? 'Later' : ''}</span>
            <button
              onClick={() => handleToggleRead(article)}
              className="ml-4 px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800 hover:bg-blue-200"
            >
              {article.is_read ? 'Mark as Unread' : 'Mark as Read'}
            </button>
            <button
              onClick={() => handleToggleLater(article)}
              className="ml-2 px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800 hover:bg-green-200"
            >
              {article.is_later ? 'Remove from Later' : 'Add to Later'}
            </button>
          </div>
        </li>
      ))}
    </ul>
  );
};

export default ArticleList;