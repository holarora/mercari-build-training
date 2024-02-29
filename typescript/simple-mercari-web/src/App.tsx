import { useState } from 'react';
import './App.css';
import { ItemList } from './components/ItemList';
import { ItemSearch, Listing } from './components/Listing';

function App() {
  // reload ItemList after Listing complete
  const [reload, setReload] = useState(true);
  const [reloadKeyword, setReloadKeyword] = useState(false);
  const [keyword, setKeyword] = useState("");
  return (
    <div>
      <header className='Title'>
        <p>
          <b>Simple Mercari</b>
        </p>
      </header>
      <div>
        <Listing onListingCompleted={() => setReload(true)} />
      </div>
      <div>
        <ItemSearch onListingCompleted={(keyword) => {
          setReloadKeyword(true);
          setKeyword(keyword);
        }} onClick={() => {setReload(true)}} />
      </div>
      <div>
        <ItemList
          keyword={keyword}
          reload={reload}
          reloadKeyword={reloadKeyword}
          onLoadCompleted={() => {
            setReload(false);
            setReloadKeyword(false);
          }}
        />
      </div>
    </div>
  )
}

export default App;