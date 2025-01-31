// @Libs
import React, { useCallback, useEffect, useMemo } from 'react';
import { generatePath, useHistory } from 'react-router-dom';
import { Button, Dropdown, message, Modal } from 'antd';
import { snakeCase } from 'lodash';
// @Components
import { DropDownList } from '@molecule/DropDownList';
import { ListItem } from '@molecule/ListItem';
import { ListItemDescription } from '@atom/ListItemDescription';
// @Icons
import PlusOutlined from '@ant-design/icons/lib/icons/PlusOutlined';
import DeleteOutlined from '@ant-design/icons/lib/icons/DeleteOutlined';
import CodeOutlined from '@ant-design/icons/lib/icons/CodeOutlined';
import EditOutlined from '@ant-design/icons/lib/icons/EditOutlined';
// @Services
import ApplicationServices from '@service/ApplicationServices';
// @Types
import { SourceConnector } from '@catalog/sources/types';
import { CommonSourcePageProps } from '@page/SourcesPage';
import { withHome } from '@molecule/Breadcrumbs/Breadcrumbs.types';
// @Styles
import styles from './SourcesList.module.less';
// @Sources
import { allSources } from '@catalog/sources/lib';
// @Routes
import { sourcesPageRoutes } from '@page/SourcesPage/SourcesPage.routes';
import { DropDownListItem } from '@molecule/DropDownList';
// @Utils
import { sourcePageUtils } from '@page/SourcesPage/SourcePage.utils';
import { taskLogsPageRoute } from '@page/TaskLogs/TaskLogsPage';
import { withProgressBar } from '@./lib/components/components';
import ExclamationCircleOutlined from '@ant-design/icons/lib/icons/ExclamationCircleOutlined';

const SourcesList = ({ projectId, sources, updateSources, setBreadcrumbs }: CommonSourcePageProps) => {
  const history = useHistory();

  const services = useMemo(() => ApplicationServices.get(), []);

  const sourcesMap = useMemo<{ [key: string]: SourceConnector }>(
    () =>
      allSources.reduce(
        (accumulator: { [key: string]: SourceConnector }, current: SourceConnector) => ({
          ...accumulator,
          [snakeCase(current.id)]: current
        }),
        {}
      ),
    []
  );

  const isFirstSingerType = useCallback((list: DropDownListItem[], item: DropDownListItem, index: number) => {
    return !item.isSingerType && list[index + 1]?.isSingerType ? styles.isFirstSingerTap : undefined;
  }, []);

  const dropDownList = useMemo(() => <DropDownList
    list={allSources.map((src: SourceConnector) => ({
      title: src.displayName,
      id: src.id,
      icon: src.pic,
      link: generatePath(sourcesPageRoutes.addExact, { source: src.id }),
      isSingerType: src.isSingerType
    }))}
    getClassName={isFirstSingerType}
    filterPlaceholder="Filter by source name or id"
  />, [isFirstSingerType]);

  const handleDeleteAction = useCallback(
    (sourceId: string) => () => {
      const updatedSources = [...sources.filter((source: SourceData) => sourceId !== source.sourceId)];

      services.storageService.save('sources', { sources: updatedSources }, projectId).then(() => {
        updateSources({ sources: updatedSources });

        message.success('Sources list successfully updated');
      });
    },
    [sources, updateSources, services.storageService, projectId]
  );

  const handleEditAction = useCallback((id: string) => () => history.push(generatePath(sourcesPageRoutes.editExact, { sourceId: id })), [history]);

  useEffect(() => {
    setBreadcrumbs(withHome({
      elements: [
        { title: 'Sources', link: sourcesPageRoutes.root },
        {
          title: 'Sources List'
        }
      ]
    }));
  }, [setBreadcrumbs]);

  if (sources.length === 0) {
    return (
      <div className={styles.empty}>
        <h3 className="text-2xl">Sources list is still empty</h3>
        <div>
          <Dropdown placement="bottomCenter" trigger={['click']} overlay={dropDownList}>
            <Button type="primary" size="large" icon={<PlusOutlined />}>Add source</Button>
          </Dropdown>
        </div>
      </div>
    );
  }

  return (
    <>
      <div className="mb-5">
        <Dropdown trigger={['click']} overlay={dropDownList}>
          <Button type="primary" icon={<PlusOutlined />}>Add source</Button>
        </Dropdown>
      </div>

      <ul>
        {sources.map((src: SourceData) => {
          const reference = sourcesMap[src.sourceProtoType];

          return <ListItem
            description={<ListItemDescription render={reference.displayName} />}
            title={sourcePageUtils.getTitle(src)}
            icon={reference?.pic}
            id={src.sourceId}
            key={src.sourceId}
            actions={[
              { onClick: () => {
                withProgressBar({
                  estimatedMs: 1000,
                  callback: async() => {
                    for (let i = 0; i < src.collections.length; i++) {
                      await services.backendApiClient.post('/tasks', undefined, { proxy: true, urlParams: {
                        source: `${services.activeProject.id}.${src.sourceId}`, collection: src.collections[i].name,
                        project_id: services.activeProject.id }
                      });
                    }
                    history.push(generatePath(taskLogsPageRoute, { sourceId: src.sourceId }));
                  }
                })

              }, title: 'Schedule All Tasks', icon: <CodeOutlined /> },
              { onClick: () => history.push(generatePath(taskLogsPageRoute, { sourceId: src.sourceId })), title: 'View logs', icon: <CodeOutlined /> },
              { onClick: () => history.push(generatePath(sourcesPageRoutes.editExact, { sourceId: src.sourceId })), title: 'Edit', icon: <EditOutlined /> },
              { onClick: () => {
                const updatedSources = [...sources.filter((source: SourceData) => src.sourceId !== source.sourceId)];

                services.storageService.save('sources', { sources: updatedSources }, projectId).then(() => {
                  updateSources({ sources: updatedSources });

                  message.success('Sources list successfully updated');
                });

              }, title: 'Delete', icon: <DeleteOutlined /> }
            ]}
          />
        })}
      </ul>
    </>
  );
};

SourcesList.displayName = 'SourcesList';

export { SourcesList };
