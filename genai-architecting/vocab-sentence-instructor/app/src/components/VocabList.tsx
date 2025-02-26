import React from 'react';
import styled from 'styled-components';
import { VocabItem } from '../types/vocab';

interface VocabListProps {
  vocabs: VocabItem[];
  loading: boolean;
}

const VocabListContainer = styled.div`
  margin-top: 2rem;
  width: 100%;
`;

const LoadingMessage = styled.div`
  text-align: center;
  font-size: 1.2rem;
  margin: 2rem 0;
  color: #666;
`;

const VocabCard = styled.div`
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  margin-bottom: 1rem;
  transition: transform 0.2s ease;

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }
`;

const VocabHeader = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
`;

const EnglishWord = styled.h3`
  margin: 0;
  color: #333;
  font-size: 1.4rem;
`;

const GermanWord = styled.h3`
  margin: 0;
  color: #1a73e8;
  font-size: 1.4rem;
  font-weight: 600;
`;

const PartsContainer = styled.div`
  margin-top: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px solid #eee;
`;

const PartLabel = styled.span`
  font-size: 0.9rem;
  color: #666;
  margin-right: 0.5rem;
`;

const PartsList = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 0.25rem;
`;

const PartTag = styled.span`
  background-color: #e8f0fe;
  color: #1a73e8;
  padding: 0.25rem 0.75rem;
  border-radius: 16px;
  font-size: 0.9rem;
`;

const NoVocabsMessage = styled.div`
  text-align: center;
  padding: 2rem;
  color: #666;
  font-size: 1.1rem;
`;

const VocabList: React.FC<VocabListProps> = ({ vocabs, loading }) => {
  if (loading) {
    return <LoadingMessage>Loading vocabulary...</LoadingMessage>;
  }

  if (vocabs.length === 0) {
    return <NoVocabsMessage>No vocabulary items to display. Enter a topic to get started!</NoVocabsMessage>;
  }

  return (
    <VocabListContainer>
      {vocabs.map((vocab, index) => (
        <VocabCard key={index}>
          <VocabHeader>
            <EnglishWord>{vocab.english}</EnglishWord>
            <GermanWord>{vocab.german}</GermanWord>
          </VocabHeader>
          {vocab.parts && vocab.parts.length > 0 && (
            <PartsContainer>
              <PartLabel>Word parts:</PartLabel>
              <PartsList>
                {vocab.parts.map((part, partIndex) => (
                  <PartTag key={partIndex}>{part}</PartTag>
                ))}
              </PartsList>
            </PartsContainer>
          )}
        </VocabCard>
      ))}
    </VocabListContainer>
  );
};

export default VocabList;
