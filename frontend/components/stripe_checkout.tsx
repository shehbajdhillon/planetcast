import { useMutation, gql } from '@apollo/client';
import { useStripe } from '@stripe/react-stripe-js';
import { Spacer, Stack } from '@chakra-ui/react';
import Button from './button';
import { Button as ChakraButton } from '@chakra-ui/react';

const CREATE_STRIPE_CHECKOUT = gql`
  mutation CreateStripeCheckout($teamSlug: String!, $lineItems:[LineItemInput!]!) {
    createCheckoutSession(teamSlug: $teamSlug, lineItems: $lineItems) {
      sessionId
    }
  }
`;

interface StripeCheckoutFormProps {
  teamSlug: string;
  subscriptionActive?: boolean;
};

const StripeCheckoutForm: React.FC<StripeCheckoutFormProps> = ({ teamSlug, subscriptionActive }) => {

  const [createCheckoutSession, { loading, error }] = useMutation(CREATE_STRIPE_CHECKOUT);

  const stripe = useStripe();

  const handleCheckout = async () => {
    const lineItems = [
      {
        priceData: {
          currency: 'usd',
          unitAmount: 1500,
          productName: 'Coffee Mug'
        },
        quantity: 2
      }
    ];

    const response = await createCheckoutSession({ variables: {
      teamSlug,
      lineItems,
    } });
    const sessionId = response.data?.createCheckoutSession?.sessionId;

    const res = await stripe?.redirectToCheckout({ sessionId });

    if (res?.error) {
      console.log('[stripe error]', res.error);
    }

  };

  return (
    <Stack w="full" direction={{ base: "column", md: "row" }} overflow={"auto"}>

      <Button isDisabled={loading} flip>
        Update Subscription
      </Button>
      <Button isDisabled={loading} flip hidden={!subscriptionActive}>
        Buy Minutes
      </Button>
      <Button flip>
        Manage Invoices
      </Button>

      <Spacer />

      <ChakraButton colorScheme='red' hidden={!subscriptionActive}>
        Cancel Subscription
      </ChakraButton>

    </Stack>
  );
};

export default StripeCheckoutForm;
