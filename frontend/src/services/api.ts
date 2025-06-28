/**
 * @fileoverview フィードアプリケーション用のAxios APIクライアント設定
 * 
 * このファイルはバックエンドAPIとの通信に使用されるAxiosインスタンスを提供します。
 * 環境変数を使用してAPIベースURLを設定し、全てのリクエストに共通ヘッダーを適用します。
 * 
 * @author FeedApp Team
 * @version 1.0.0
 */

import axios from 'axios';

/**
 * APIのベースURL
 * 環境変数NEXT_PUBLIC_API_BASE_URLが設定されている場合はその値を使用し、
 * 未設定の場合はローカル開発用のデフォルトURLを使用
 */
const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1';

/**
 * 設定済みのAxiosインスタンス
 * 
 * バックエンドAPIとの通信に使用するAxiosクライアントです。
 * ベースURLと共通ヘッダーが事前に設定されており、
 * 全てのAPIリクエストでJSON形式の通信を行います。
 * 
 * @example
 * ```typescript
 * import api from './api';
 * 
 * // GET リクエストの例
 * const response = await api.get('/folders');
 * 
 * // POST リクエストの例
 * const newFolder = await api.post('/folders', { name: 'New Folder' });
 * ```
 */
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export default api;
