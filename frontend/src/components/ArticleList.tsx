/**
 * @fileoverview 記事一覧表示コンポーネント
 * 
 * このファイルは記事の一覧を表示し、ユーザーが記事の読了状態や
 * 「後で見る」状態を操作できるインタラクティブなUIコンポーネントを提供します。
 * 
 * @author FeedApp Team
 * @version 1.0.0
 */

import React from 'react';
import { Article, articleService } from '../services/article';
import { useSWRConfig } from 'swr';

/**
 * ArticleListコンポーネントのProps型定義
 * 
 * @interface ArticleListProps
 */
interface ArticleListProps {
  /** 表示する記事の配列 */
  articles: Article[];
}

/**
 * 記事一覧表示コンポーネント
 * 
 * 記事の配列を受け取り、各記事を読みやすいカード形式で表示します。
 * 各記事には以下の機能があります：
 * - 記事タイトルのクリックで外部リンクを開く
 * - 既読/未読の切り替え
 * - 「後で見る」の追加/削除
 * - 記事の公開日表示
 * - 記事本文のプレビュー（150文字まで）
 * 
 * @param {ArticleListProps} props - コンポーネントのプロパティ
 * @returns {JSX.Element} 記事一覧のJSX要素
 * 
 * @example
 * ```tsx
 * import ArticleList from './ArticleList';
 * 
 * const articles = [
 *   {
 *     id: '1',
 *     title: 'サンプル記事',
 *     content: '記事の内容...',
 *     url: 'https://example.com/article1',
 *     is_read: false,
 *     is_later: false,
 *     published_at: '2023-01-01',
 *     // ... その他のプロパティ
 *   }
 * ];
 * 
 * <ArticleList articles={articles} />
 * ```
 */
const ArticleList: React.FC<ArticleListProps> = ({ articles }) => {
  const { mutate } = useSWRConfig();
  
  /**
   * 記事の既読/未読状態を切り替えます
   * 
   * APIリクエストを送信して記事の既読状態を更新し、
   * 成功時にはSWRキャッシュを無効化して画面を更新します。
   * エラー時にはコンソールログとアラートでユーザーに通知します。
   * 
   * @param {Article} article - 状態を変更する記事オブジェクト
   * @returns {Promise<void>} 非同期処理のPromise
   */
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

  /**
   * 記事の「後で見る」状態を切り替えます
   * 
   * APIリクエストを送信して記事の「後で見る」状態を更新し、
   * 成功時にはSWRキャッシュを無効化して画面を更新します。
   * エラー時にはコンソールログとアラートでユーザーに通知します。
   * 
   * @param {Article} article - 状態を変更する記事オブジェクト
   * @returns {Promise<void>} 非同期処理のPromise
   */
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