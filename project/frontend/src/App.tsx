import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { TopicList } from './components/TopicList';
import { TopicDetail } from './components/TopicDetail';
import { SearchBar } from './components/SearchBar';
import { SearchResults } from './components/SearchResults';
import './App.css';

function App() {
  return (
    <BrowserRouter>
      <div className="app">
        <header className="app-header">
          <h1>ReSQL Forum Dashboard</h1>
          <SearchBar />
        </header>
        <main className="app-main">
          <Routes>
            <Route path="/" element={<TopicList />} />
            <Route path="/topics/:topicId" element={<TopicDetail />} />
            <Route path="/search" element={<SearchResults />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;
