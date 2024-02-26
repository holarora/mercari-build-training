import React, { useState } from 'react';

interface formDataType {
  name: string;
  category: string;
  image: string;
};

const server = process.env.REACT_APP_API_URL || 'http://127.0.0.1:9000';

interface Prop {
    onKeywordCompleted?: () => void;
}

export const ItemSearch: React.FC<Prop> = (props) => {
    const { onKeywordCompleted } = props;
    const initialState = {
      name: "",
      category: "",
      image: "",
    };
    const [keyword, setKeyword] = useState<string>("");

    const onKeywordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setKeyword(event.target.value);
      };
    const onSubmit = (event: React.FormEvent<HTMLFormElement>) => {
      event.preventDefault();

      fetch(server.concat('/search?keyword=' + encodeURIComponent(keyword)), {
        method: 'GET',
        mode: 'cors',
      })
        .then(response => {
          console.log('GET status:', response.statusText);
          onKeywordCompleted && onKeywordCompleted();
        })
        .catch((error) => {
          console.error('GET error:', error);
        })
    };
    return (
      <div className='Listing'>
        <form onSubmit={onSubmit}>
          <div>
            <input type='text' name='keyword' id='keyword' placeholder='key word' onChange={onKeywordChange} required className='TextInput'/>
            <button type='submit' className='SubmitInput'>Search</button>
          </div>
        </form>
      </div>
    );
};
