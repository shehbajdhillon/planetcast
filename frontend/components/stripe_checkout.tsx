import { useMutation, gql } from '@apollo/client';
import { useStripe } from '@stripe/react-stripe-js';
import { Spacer, Stack } from '@chakra-ui/react';
import Button from './button';
import { Button as ChakraButton } from '@chakra-ui/react';
import { useRouter } from 'next/router';

const CREATE_STRIPE_CHECKOUT = gql`
  mutation CreateStripeCheckout($teamSlug: String!, $lineItems:[LineItemInput!]!) {
    createCheckoutSession(teamSlug: $teamSlug, lineItems: $lineItems) {
      sessionId
    }
  }
`;

const CREATE_STRIPE_PORTAL = gql`
  mutation CreateStripePortalSession($teamSlug: String!) {
    createPortalSession(teamSlug: $teamSlug) {
      sessionUrl
    }
  }
`;

interface StripeCheckoutFormProps {
  teamSlug: string;
  subscriptionActive?: boolean;
};

const StripeCheckoutForm: React.FC<StripeCheckoutFormProps> = ({ teamSlug, subscriptionActive }) => {

  const router = useRouter();

  const [createCheckoutSession, { loading, error }] = useMutation(CREATE_STRIPE_CHECKOUT);
  const [createPortalSession, { loading: portalLoading, error: portalError }] = useMutation(CREATE_STRIPE_PORTAL);

  const stripe = useStripe();

  const handleCheckout = async (lookUpKey: string) => {
    const response = await createCheckoutSession({ variables: { teamSlug, lookUpKey } });
    const sessionId = response.data?.createCheckoutSession?.sessionId;
    const res = await stripe?.redirectToCheckout({ sessionId });
    if (res?.error) {
      console.log('[stripe error]', res.error.message);
    }
  };

  const handlePortalSession = async () => {
    const response = await createPortalSession({ variables: { teamSlug } });
    const sessionUrl = response.data?.createPortalSession?.sessionUrl;
    if (sessionUrl) {
      router.push(sessionUrl);
    } else {
      console.log('[stripe error]', portalError?.message);
    }
  };

  return (
    <Stack w="full" direction={{ base: "column", md: "row" }} overflow={"auto"}>

      <Button isDisabled={loading} flip={true}>
        Update Subscription
      </Button>
      <Button isDisabled={loading} flip={true} hidden={!subscriptionActive}>
        Buy Minutes
      </Button>
      <Button isDisabled={portalLoading} onClick={handlePortalSession} flip={true}>
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
