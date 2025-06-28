-- 20250628081751_insert_test_data.sql

-- foldersテーブルにテストデータを挿入
INSERT INTO folders (id, name, user_id) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '未分類', NULL),
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'ニュース', NULL),
('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', '技術ブログ', NULL)
ON CONFLICT (id) DO NOTHING;

-- feedsテーブルにテストデータを挿入
INSERT INTO feeds (id, name, url, plugin_type, folder_id, update_interval) VALUES
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'Google News', 'https://news.google.com/rss?hl=ja&gl=JP&ceid=JP:ja', 'rss', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 60),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'Qiita', 'https://qiita.com/popular-items/feed', 'rss', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 120)
ON CONFLICT (id) DO NOTHING;

-- articlesテーブルにテストデータを挿入 (一部のみ)
INSERT INTO articles (id, feed_id, title, content, url, published_at, is_read, is_later) VALUES
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'テスト記事1', 'これはテスト記事1の内容です。', 'http://example.com/article1', '2025-06-28 10:00:00+09', FALSE, FALSE),
('10eebc99-9c0b-4ef8-bb6d-6bb9bd380a17', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'テスト記事2', 'これはテスト記事2の内容です。', 'http://example.com/article2', '2025-06-28 11:00:00+09', FALSE, TRUE)
ON CONFLICT (id) DO NOTHING;

-- pluginsテーブルにテストデータを挿入
INSERT INTO plugins (id, name, file_path, enabled) VALUES
('20eebc99-9c0b-4ef8-bb6d-6bb9bd380a18', 'rss_plugin', '/path/to/rss_plugin.so', TRUE),
('30eebc99-9c0b-4ef8-bb6d-6bb9bd380a19', 'custom_plugin_example', '/path/to/custom_plugin.so', TRUE)
ON CONFLICT (id) DO NOTHING;
