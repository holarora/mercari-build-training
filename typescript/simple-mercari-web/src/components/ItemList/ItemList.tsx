import React, { useEffect, useState } from 'react';

interface Item {
  id: number;
  name: string;
  category: string;
  image_name: string;
};

const server = process.env.REACT_APP_API_URL || 'http://127.0.0.1:9000';
const placeholderImage = process.env.PUBLIC_URL + '/logo192.png';

interface Prop {
  keyword: string;
  reload?: boolean;
  reloadKeyword?: boolean;
  onLoadCompleted?: () => void;
}

export const ItemList: React.FC<Prop> = (props) => {
  const { keyword, reload = true, reloadKeyword = false, onLoadCompleted } = props;
  const [items, setItems] = useState<Item[]>([])
  const fetchItems = () => {
    fetch(server.concat('/items'),
      {
        method: 'GET',
        mode: 'cors',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
      })
      .then(response => response.json())
      .then(data => {
        console.log('GET success:', data);
        setItems(data.items);
        onLoadCompleted && onLoadCompleted();
      })
      .catch(error => {
        console.error('GET error:', error)
      })
  }
  const fetchItemsByKeyword = () => {
    fetch(server.concat('/search?keyword=', encodeURIComponent(keyword)),
      {
        method: 'GET',
        mode: 'cors',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
      })
      .then(response => response.json())
      .then(data => {
        console.log('GET success:', data);
        setItems(data.items);
        onLoadCompleted && onLoadCompleted();
      })
      .catch(error => {
        console.error('GET error:', error)
      })
  }

  useEffect(() => {
    if (reload) {
      fetchItems();
    } else if (reloadKeyword) {
      fetchItemsByKeyword();
    }
  });

  return (
    <div className='Wrapper'>
      {items.map((item) => {
        return (
          <div key={item.id} className='ItemList'>
            <img
              src={"http://localhost:9000/image/" + item.image_name}
              alt={placeholderImage}
              className='Image' />
            <p>
              <span className='Item'>Name: {item.name}</span>
              <br />
              <span className='Item'>Category: {item.category}</span>
            </p>
          </div>
        )
      })}
    </div>
  )
};
