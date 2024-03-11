import React, { useEffect, useState } from 'react';

interface Item {
  id: number;
  name: string;
  category: string;
  image_name: string;
};

const server = process.env.REACT_APP_API_URL || 'http://localhost:9000';
const placeholderImage = process.env.PUBLIC_URL + '/logo192.png';

interface Prop {
  reload?: boolean;
  onLoadCompleted?: () => void;
}

export const ItemList: React.FC<Prop> = (props) => {
  const { reload = true, onLoadCompleted } = props;
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
        // setItems(data);
        onLoadCompleted && onLoadCompleted();
      })
      .catch(error => {
        console.error('GET error:', error)
      })
  }

  useEffect(() => {
    if (reload) {
      fetchItems();
    }
    //useEffectの第二引数(依存関係配列)にfetchItemsを入れることで、この関数を変更されたときに、fetchItemsで取得データの変更などが
    //コンポーネントに反映されるイメージ。
  }, [reload]);

  const Like = () => {
    const [count, setCout] = useState(0);
    const handleClick = () => {
      setCout(count+1);
    };
    return(
      <span className='likeButton' onClick={handleClick}>♡ {count}</span>
    );
  }

  return (
    <div className='Container'>
      {/* console.log('WHAT items', items) */}
      {/* console.log(items) */}
      {items.map((item) => {
      const imageUrl = `${server}/image/${item.image_name}`;
        return (
          <div key={item.id} className='ItemList'>
            <div className='item_box'>
              <div className='item_box_tape'>
                 {/* TODO: Task 1: Replace the placeholder image with the item image */}
                <img src={imageUrl}  alt={placeholderImage} width="150" height="150" />
                <p>
                  <span>Name: {item.name}</span>
                  <br />
                  <span>Category: {item.category}</span>
                  <br />
                  <br />
                  <span className='likebutton'><Like /></span>
                </p>
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
};
