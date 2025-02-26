import React, { useState } from 'react';
import styled from 'styled-components';

interface TopicSearchProps {
  onSearch: (topic: string) => void;
  loading: boolean;
}

const SearchContainer = styled.div`
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
`;

const SearchForm = styled.form`
  display: flex;
  gap: 0.5rem;
  width: 100%;
`;

const SearchInput = styled.input`
  flex: 1;
  padding: 0.75rem 1rem;
  font-size: 1rem;
  border: 2px solid #ddd;
  border-radius: 8px;
  outline: none;
  transition: border-color 0.2s;

  &:focus {
    border-color: #1a73e8;
  }
`;

const SearchButton = styled.button`
  background-color: #1a73e8;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0 1.5rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #1557b0;
  }

  &:disabled {
    background-color: #9fc1ee;
    cursor: not-allowed;
  }
`;

const TopicSearch: React.FC<TopicSearchProps> = ({ onSearch, loading }) => {
  const [topic, setTopic] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (topic.trim()) {
      onSearch(topic.trim());
    }
  };

  return (
    <SearchContainer>
      <SearchForm onSubmit={handleSubmit}>
        <SearchInput
          type="text"
          placeholder="Enter a topic in English (e.g., food, travel, technology)"
          value={topic}
          onChange={(e) => setTopic(e.target.value)}
          disabled={loading}
        />
        <SearchButton type="submit" disabled={loading || !topic.trim()}>
          {loading ? 'Loading...' : 'Search'}
        </SearchButton>
      </SearchForm>
    </SearchContainer>
  );
};

export default TopicSearch;
