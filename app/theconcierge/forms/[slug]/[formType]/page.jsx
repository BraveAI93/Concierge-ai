'use client';
import FormPage from '@/components/FormPage';

export default function FormPageRoute({ params }) {
  return <FormPage profileSlug={params.slug} formTypeSlug={params.formType} />;
}
