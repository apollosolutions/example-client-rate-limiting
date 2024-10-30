const resolvers = {
  Query: {
    hello: () => {
      return "Hello World!";
    },
    helloExpensive: () => {
      return "Hello from a really big world!";
    },
  },
};

export default resolvers;
