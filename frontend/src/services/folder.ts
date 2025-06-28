/**
 * @fileoverview フォルダ関連のAPIサービスと型定義
 * 
 * このファイルはフォルダの作成、取得、更新、削除に関する
 * APIクライアント機能を提供します。TypeScriptの型安全性を
 * 活用してフロントエンドとバックエンド間のデータ整合性を保証します。
 * 
 * @author FeedApp Team
 * @version 1.0.0
 */

import api from './api';

/**
 * フォルダのデータ型定義
 * 
 * バックエンドのFolder構造体と対応するTypeScriptインターフェースです。
 * フィードを分類・整理するためのコンテナとして機能します。
 * 
 * @interface Folder
 */
export interface Folder {
  /** フォルダの一意識別子（UUID形式） */
  id: string;
  /** フォルダの表示名 */
  name: string;
  /** フォルダの作成日時（ISO 8601形式） */
  created_at: string;
}

/**
 * フォルダAPIサービス
 * 
 * フォルダのCRUD操作を行うためのAPIクライアント機能を提供します。
 * 全ての関数は非同期処理でPromiseを返し、エラーハンドリングは
 * 呼び出し元で行うことを想定しています。
 * 
 * @namespace folderService
 * @example
 * ```typescript
 * import { folderService } from './folder';
 * 
 * // 全フォルダ取得
 * const folders = await folderService.getAllFolders();
 * 
 * // 新規フォルダ作成
 * const newFolder = await folderService.createFolder('技術記事');
 * 
 * // フォルダ更新
 * const updated = await folderService.updateFolder(folder.id, '技術ブログ');
 * ```
 */
export const folderService = {
  /**
   * 全てのフォルダを取得します
   * 
   * データベースに保存されている全てのフォルダを配列で取得します。
   * 作成日時順でソートされて返されます。
   * 
   * @returns {Promise<Folder[]>} フォルダ配列のPromise
   * @throws {Error} APIリクエストが失敗した場合
   * 
   * @example
   * ```typescript
   * try {
   *   const folders = await folderService.getAllFolders();
   *   console.log(`取得したフォルダ数: ${folders.length}`);
   * } catch (error) {
   *   console.error('フォルダ取得エラー:', error);
   * }
   * ```
   */
  getAllFolders: async (): Promise<Folder[]> => {
    const response = await api.get<Folder[]>('/folders');
    return response.data;
  },

  /**
   * 指定されたIDのフォルダを取得します
   * 
   * @param {string} id - フォルダのID（UUID形式）
   * @returns {Promise<Folder>} フォルダオブジェクトのPromise
   * @throws {Error} 指定されたIDのフォルダが存在しない場合（404）
   * @throws {Error} APIリクエストが失敗した場合
   * 
   * @example
   * ```typescript
   * try {
   *   const folder = await folderService.getFolderById('550e8400-e29b-41d4-a716-446655440000');
   *   console.log(`フォルダ名: ${folder.name}`);
   * } catch (error) {
   *   console.error('フォルダ取得エラー:', error);
   * }
   * ```
   */
  getFolderById: async (id: string): Promise<Folder> => {
    const response = await api.get<Folder>(`/folders/${id}`);
    return response.data;
  },

  /**
   * 新しいフォルダを作成します
   * 
   * 指定された名前で新しいフォルダを作成します。
   * IDと作成日時は自動的に設定されます。
   * 
   * @param {string} name - フォルダ名（必須、空文字列不可）
   * @returns {Promise<Folder>} 作成されたフォルダオブジェクトのPromise
   * @throws {Error} フォルダ名が不正な場合（400）
   * @throws {Error} APIリクエストが失敗した場合
   * 
   * @example
   * ```typescript
   * try {
   *   const newFolder = await folderService.createFolder('ニュース');
   *   console.log(`作成されたフォルダID: ${newFolder.id}`);
   * } catch (error) {
   *   console.error('フォルダ作成エラー:', error);
   * }
   * ```
   */
  createFolder: async (name: string): Promise<Folder> => {
    const response = await api.post<Folder>('/folders', { name });
    return response.data;
  },

  /**
   * 既存のフォルダ情報を更新します
   * 
   * 指定されたIDのフォルダの名前を更新します。
   * 
   * @param {string} id - 更新するフォルダのID（UUID形式）
   * @param {string} name - 新しいフォルダ名（必須、空文字列不可）
   * @returns {Promise<Folder>} 更新されたフォルダオブジェクトのPromise
   * @throws {Error} 指定されたIDのフォルダが存在しない場合（404）
   * @throws {Error} フォルダ名が不正な場合（400）
   * @throws {Error} APIリクエストが失敗した場合
   * 
   * @example
   * ```typescript
   * try {
   *   const updated = await folderService.updateFolder(folderId, '技術ブログ');
   *   console.log(`更新後のフォルダ名: ${updated.name}`);
   * } catch (error) {
   *   console.error('フォルダ更新エラー:', error);
   * }
   * ```
   */
  updateFolder: async (id: string, name: string): Promise<Folder> => {
    const response = await api.put<Folder>(`/folders/${id}`, { name });
    return response.data;
  },

  /**
   * 指定されたIDのフォルダを削除します
   * 
   * フォルダを完全に削除します。この操作は取り消すことができません。
   * フォルダに関連付けられたフィードがある場合の動作は
   * バックエンドの実装に依存します。
   * 
   * @param {string} id - 削除するフォルダのID（UUID形式）
   * @returns {Promise<void>} 削除完了のPromise（戻り値なし）
   * @throws {Error} 指定されたIDのフォルダが存在しない場合（404）
   * @throws {Error} APIリクエストが失敗した場合
   * 
   * @example
   * ```typescript
   * try {
   *   await folderService.deleteFolder(folderId);
   *   console.log('フォルダが削除されました');
   * } catch (error) {
   *   console.error('フォルダ削除エラー:', error);
   * }
   * ```
   */
  deleteFolder: async (id: string): Promise<void> => {
    await api.delete(`/folders/${id}`);
  },
};
