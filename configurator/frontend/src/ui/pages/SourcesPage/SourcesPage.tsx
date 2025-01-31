// @Libs
import React, { Dispatch, SetStateAction, useMemo } from 'react';
import { Route, Switch } from 'react-router-dom';
// @Routes
import { sourcesPageRoutes } from './SourcesPage.routes';
// @Components
import { SourcesList } from './partials/SourcesList';
import { SourceEditor } from './partials/SourceEditor';
import { CenteredError, CenteredSpin } from '@./lib/components/components';
// @Services
import ApplicationServices from '@service/ApplicationServices';
// @Styles
import './SourcesPage.less';
// @Hocs
import { getComponent } from '@hocs/getComponent';
// @Types
import { BreadcrumbsProps } from '@molecule/Breadcrumbs/Breadcrumbs.types';
import { PageProps } from '@./navigation';
// @Hooks
import useLoader from '@hooks/useLoader';

export interface CollectionSourceData {
  sources: SourceData[];
  _lastUpdated?: string;
}

export interface CommonSourcePageProps {
  sources: SourceData[];
  projectId: string;
  updateSources: Dispatch<SetStateAction<CollectionSourceData>>;
  setBreadcrumbs: (breadcrumbs: BreadcrumbsProps) => void;
  editorMode?: 'edit' | 'add';
}

const SourcesPage = (props: PageProps) => {
  const services = useMemo(() => ApplicationServices.get(), []);

  const [error, sources, updateSources] = useLoader(
    async() => await services.storageService.get('sources', services.activeProject.id)
  );

  const additionalProps = useMemo(() => ({
    projectId: services.activeProject.id,
    sources: sources?.sources ?? [],
    updateSources,
    setBreadcrumbs: props.setBreadcrumbs
  }), [props.setBreadcrumbs, sources?.sources, services.activeProject.id, updateSources]);

  if (error) {
    return <CenteredError error={error} />;
  } else if (!sources) {
    return <CenteredSpin />;
  }

  return (
    <Switch>
      <Route
        path={sourcesPageRoutes.root}
        exact
        render={getComponent<CommonSourcePageProps>(SourcesList, additionalProps)}
      />
      <Route
        path={[sourcesPageRoutes.add, sourcesPageRoutes.addExact]}
        strict={false}
        exact
        render={getComponent<CommonSourcePageProps>(SourceEditor, { ...additionalProps, editorMode: 'add' })}
      />
      <Route
        path={[sourcesPageRoutes.edit, sourcesPageRoutes.editExact]}
        strict={false}
        exact
        render={getComponent<CommonSourcePageProps>(SourceEditor, { ...additionalProps, editorMode: 'edit' })}
      />
    </Switch>
  );
};

SourcesPage.displayName = 'SourcesPage';

export { SourcesPage };
