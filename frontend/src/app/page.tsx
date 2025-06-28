"use client";

import useSWR from 'swr';
import { Folder } from '../services/folder';
import { Article } from '../services/article';
import ArticleList from '../components/ArticleList';

export default function Home() {
  const { data: folders, error: foldersError } = useSWR<Folder[]>('/folders');
  const { data: articles, error: articlesError } = useSWR<Article[]>('/articles');

  if (foldersError || articlesError) return <div>Failed to load data</div>;
  if (!folders || !articles) return <div>Loading...</div>;

  return (
    <div className="flex">
      <div className="w-1/4 p-4 border-r">
        <h1 className="text-2xl font-bold mb-4">Folders</h1>
        <ul>
          {folders.map((folder) => (
            <li key={folder.id} className="mb-2 p-2 border rounded">
              {folder.name}
            </li>
          ))}
        </ul>
      </div>
      <div className="w-3/4 p-4">
        <h1 className="text-2xl font-bold mb-4">Articles</h1>
        <ArticleList articles={articles} />
      </div>
    </div>
  );
}