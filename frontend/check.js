const checkBookmarks = async () => {
  const response = await fetch('/api/public/bookmarks');
  const bookmarks = await response.json();
  
  console.log('Total bookmarks:', bookmarks.length);
  console.log('\nFirst 5 bookmarks:');
  bookmarks.slice(0, 5).forEach((b, i) => {
    console.log(`${i+1}. ID: ${b.id}, Title: ${b.title}, Icon: '${b.icon}'`);
  });
  
  const emptyIconCount = bookmarks.filter(b => !b.icon || b.icon === '').length;
  console.log(`\nBookmarks with empty icon: ${emptyIconCount}`);
};

checkBookmarks();
