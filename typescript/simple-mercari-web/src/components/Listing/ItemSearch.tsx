import React, { useState } from 'react';

interface Prop {
  onListingCompleted?: (keyword: string) => void;
  onClick: () => void;
}

export const ItemSearch: React.FC<Prop> = (props) => {
  const { onListingCompleted, onClick } = props;
  const [keyword, setKeyword] = useState<string>("");

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    onListingCompleted && onListingCompleted(keyword);
  };
  return (
    <div className='ItemSearch'>
      <form onSubmit={onSubmit}>
        <div>
          <input
            type='text'
            name='keyword'
            id='keyword'
            placeholder='keyword'
            value={keyword}
            onChange={(event) => setKeyword(event.target.value)}
            required
            className='TextInput'
          />
          <button type='submit' className='SubmitSearch'>Search</button>
        </div>
      </form>
      <div>
        <button onClick={onClick} className='SubmitShowAll'>Show All</button>
      </div>
    </div>
  );
}
