import React from 'react';
import styled from 'styled-components';
import { VocabItem } from '../types/vocab';

interface DownloadButtonProps {
  vocabs: VocabItem[];
  topic: string | null;
  disabled: boolean;
}

const Button = styled.button`
  background-color: #34a853;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  transition: background-color 0.2s;
  margin: 1rem auto;

  &:hover {
    background-color: #2e8b46;
  }

  &:disabled {
    background-color: #a8d5b1;
    cursor: not-allowed;
  }

  svg {
    width: 1.2rem;
    height: 1.2rem;
  }
`;

const DownloadButton: React.FC<DownloadButtonProps> = ({ vocabs, topic, disabled }) => {
  const handleDownload = () => {
    if (vocabs.length === 0 || !topic) return;

    // Create a JSON object with metadata
    const jsonData = {
      topic,
      timestamp: new Date().toISOString(),
      count: vocabs.length,
      vocabs
    };

    // Convert to JSON string
    const jsonString = JSON.stringify(jsonData, null, 2);
    
    // Create a blob with the JSON data
    const blob = new Blob([jsonString], { type: 'application/json' });
    
    // Create a URL for the blob
    const url = URL.createObjectURL(blob);
    
    // Create a temporary anchor element
    const a = document.createElement('a');
    a.href = url;
    a.download = `german-vocab-${topic.toLowerCase().replace(/\s+/g, '-')}.json`;
    
    // Trigger the download
    document.body.appendChild(a);
    a.click();
    
    // Clean up
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <Button onClick={handleDownload} disabled={disabled}>
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
      </svg>
      Download Vocabulary
    </Button>
  );
};

export default DownloadButton;
