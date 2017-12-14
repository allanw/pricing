import pandas as pd

def average_price_per_order(df):
  return df.groupby('order_id').mean().mean()['total_price']

def relative_product_popularity(df):
  top_products = df.groupby('product_id')['quantity'].sum()
  return top_products.sort_values(ascending=False)

df = pd.read_csv('orders.csv')

print('Average price per order:' , average_price_per_order(df))
print('Relative product popularity:\n', relative_product_popularity(df))
